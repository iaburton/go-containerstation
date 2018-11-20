package containerstation

import (
	"context"
	"io"
	"net/http"
	"os"
)

func (c Client) DownloadTLSCertificate(ctx context.Context, file string, perm os.FileMode) (err error) {
	var fh *os.File
	if fh, err = os.OpenFile(file, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, perm); err != nil {
		return
	}

	defer func() {
		fh.Close()
		//intentionally not capturing error from Close above
		if err != nil {
			os.Remove(file)
		}
	}()

	return c.ExportTLSCertificate(ctx, fh)
}

func (c Client) ExportTLSCertificate(ctx context.Context, w io.Writer) error {
	const apiEndpoint = `/containerstation/api/v1/tls/export/registry`

	req, err := http.NewRequest(http.MethodGet, *c.baseURL+apiEndpoint, nil)
	if err != nil {
		return err
	}

	req = req.WithContext(ctx)
	resp, err := c.hc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	//just in case, as this is a blind copy
	//5mb should be plenty for a cert/cert chain
	r := io.LimitReader(resp.Body, 1024*1024*5)
	_, err = io.Copy(w, r)
	return err
}
