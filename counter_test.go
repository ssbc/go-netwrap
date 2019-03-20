package netwrap

import (
	"io"
	"io/ioutil"
	"net"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCounter(t *testing.T) {
	r := require.New(t)

	i, o := net.Pipe()

	errc := make(chan error)
	var got []byte
	go func() {
		var err error
		got, err = ioutil.ReadAll(o)
		if err != nil {
			errc <- err
		}
		close(errc)
	}()

	ic, err := CountingWrapper(i)
	r.NoError(err)

	testStr := "128flaschenbieranderwand"
	n, err := io.Copy(ic, strings.NewReader(testStr))
	r.NoError(err)
	r.EqualValues(n, len(testStr))
	r.NoError(ic.Close())

	r.NoError(<-errc)
	r.Equal(testStr, string(got))

	cc, ok := ic.(*countingConn)
	r.True(ok)
	r.EqualValues(len(testStr), cc.tx)
	r.EqualValues(0, cc.rx)
}
