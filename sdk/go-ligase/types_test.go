package ligase

import "testing"

func TestA(t *testing.T) {
	req := &PostLoginRequest{}
	req.Address = "abc"
	t.Logf("%s", req)
}
