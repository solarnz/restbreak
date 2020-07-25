package x

import (
	"time"

	"github.com/solarnz/restbreak/src/types"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/screensaver"
	"github.com/BurntSushi/xgb/xproto"
)

type XActivityListener struct {
	XConnection   *xgb.Conn
	RootWindow    xproto.Drawable
	CheckInterval time.Duration
}

func NewXActivityListener() (*XActivityListener, error) {
	// Connect to the X Server
	X, err := xgb.NewConn()
	if err != nil {
		return nil, err
	}

	// We're going to use the screensaver extension, so we need to initialise it.
	if err = screensaver.Init(X); err != nil {
		return nil, err
	}

	// Get the root window from X via xproto
	setup := xproto.Setup(X)
	rootWindow := setup.DefaultScreen(X).Root

	listener := &XActivityListener{
		XConnection:   X,
		RootWindow:    xproto.Drawable(rootWindow),
		CheckInterval: time.Second,
	}

	return listener, nil
}

func (l *XActivityListener) Listen(stop chan interface{}, events chan types.ActivityEvent) error {

	for {
		select {
		case <-stop:
			return nil
		default:
			// Query the screensaver extension to get information on when the
			// user was last active.
			cookie := screensaver.QueryInfo(l.XConnection, l.RootWindow)
			info, err := cookie.Reply()

			if err != nil {
				return err
			}

			currentTime := time.Now()
			// Convert the value we get back from screensaver to a timestamp
			// the user was last active
			duration := -1 * time.Millisecond * time.Duration(info.MsSinceUserInput)
			timeLastActive := currentTime.Add(duration).Round(time.Second)

			events <- types.ActivityEvent{LastActiveTimestamp: timeLastActive}

			// We don't want to query the screensaver extension constantly, as
			// we don't need that high precision. I'm also not aware of any way
			// we can get a stream of events from X in order to notify us when
			// the user is active / becomes inactive.
			time.Sleep(l.CheckInterval)
		}
	}
}
