package download

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"
	"io/ioutil"

	"github.com/sasakiyudai/gopherdojo-studyroom/kadai3-2/sasakiyudai/request"
	"golang.org/x/sync/errgroup"
)

type PDownloader struct {
	url      *url.URL
	output   *os.File
	fileSize uint
	part     uint
	procs    uint
}

func newPDownloader(url *url.URL, output *os.File, fileSize uint, part uint, procs uint) *PDownloader {
	return &PDownloader{
		url:      url,
		output:   output,
		fileSize: fileSize,
		part:     part,
		procs:    procs,
	}
}

func Downloader(url *url.URL, output *os.File, fileSize uint, part uint, procs uint, isPara bool, tmpDirName string, ctx context.Context) error {
	pd := newPDownloader(url, output, fileSize, part, procs)
	if !isPara {
		fmt.Printf("%s do not accept range access: downloading by single process\n", url)
		err := pd.DownloadFile(ctx)
		if err != nil {
			return err
		}
	} else {
		grp, ctx := errgroup.WithContext(ctx)
		if err := pd.PDownload(grp, tmpDirName, procs, ctx); err != nil {
			return err
		}

		if err := grp.Wait(); err != nil {
			return err
		}
	}
	return nil
}

func (pd *PDownloader) DownloadFile(ctx context.Context) (err error) {
	resp, err := request.Request(ctx, "GET", pd.url.String(), "", "")
	if err != nil {
		return
	}
	defer func() {
		err = resp.Body.Close()
	}()

	_, err = io.Copy(pd.output, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func (pd *PDownloader) PDownload(grp *errgroup.Group, tmpDirName string, procs uint, ctx context.Context) error {
	var start, end, idx uint

	for idx = uint(0); idx < procs; idx++ {
		if idx == 0 {
			start = 0
		} else {
			start = idx*pd.part + 1
		}

		if idx == pd.procs-1 {
			end = pd.fileSize
		} else {
			end = (idx + 1) * pd.part
		}

		idx := idx

		bytes := fmt.Sprintf("bytes=%d-%d", start, end)

		grp.Go(func() error {
			fmt.Printf("grp.Go: tmpDirName: %s, bytes %s, idx: %d\n", tmpDirName, bytes, idx)
			return pd.ReqToMakeCopy(tmpDirName, bytes, idx, ctx)
		})

	}
	return nil
}

func (pd *PDownloader) ReqToMakeCopy(tmpDirName string, bytes string, idx uint, ctx context.Context) (err error) {
	resp, err := request.Request(ctx, "GET", pd.url.String(), "Range", bytes)
	if err != nil {
		return err
	}

	tmpOut, err := os.Create(tmpDirName + "/" + strconv.Itoa(int(idx)))
	if err != nil {
		return err
	}
	defer func() {
		err = tmpOut.Close()
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		if err != io.EOF && err != io.ErrUnexpectedEOF {
			return err
		}
	}

	length, err := tmpOut.Write(body)
	if err != nil {
		return err
	}
	fmt.Printf("%d/%d was downloaded len=%d\n", idx, pd.procs, length)
	return nil
}
