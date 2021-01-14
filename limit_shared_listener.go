// Writed by yijian on 2021/01/14
// 包 “golang.org/x/net/netutil” 中的 netutil.LimitListener 只针对单个 Listener，
// LimitSharedListener 可用于多个 Listener 共享连接数限制。
package lsl

import (
    "net"
    "sync"
)

func LimitSharedListener(l net.Listener, ll* ListenLimiter) net.Listener {
    return &limitSharedListener{
        Listener: l,
        limiter: ll,
    }
}

func NewListenLimiter(n int) *ListenLimiter {
    return &ListenLimiter{
        sem:      make(chan struct{}, n),
        done:     make(chan struct{}),
    }
}

type ListenLimiter struct {
    sem       chan struct{}
    closeOnce sync.Once     // ensures the done chan is only closed once
    done      chan struct{} // no values sent; closed when Close is called
}

// limitSharedListener 为接口 net.Listener 的这实现：
// Accept()
// Close()
// Addr()
type limitSharedListener struct {
    net.Listener
    limiter* ListenLimiter
}

// acquire acquires the limiting semaphore. Returns true if successfully
// accquired, false if the listener is closed and the semaphore is not
// acquired.
func (l *ListenLimiter) acquire() bool {
    select {
    case <-l.done:
        return false
    case l.sem <- struct{}{}:
        return true
    }
}
func (l *ListenLimiter) release() { <-l.sem }

func (l *limitSharedListener) Accept() (net.Conn, error) {
    acquired := l.limiter.acquire()
    // If the semaphore isn't acquired because the listener was closed, expect
    // that this call to accept won't block, but immediately return an error.
    c, err := l.Listener.Accept()
    if err != nil {
        if acquired {
            l.limiter.release()
        }
        return nil, err
    }
    return &limitListenerConn{Conn: c, release: l.limiter.release}, nil
}

func (l *limitSharedListener) Close() error {
    err := l.Listener.Close()
    l.limiter.closeOnce.Do(func() { close(l.limiter.done) })
    return err
}

type limitListenerConn struct {
    net.Conn
    releaseOnce sync.Once
    release     func()
}

func (l *limitListenerConn) Close() error {
    err := l.Conn.Close()
    l.releaseOnce.Do(l.release)
    return err
}
