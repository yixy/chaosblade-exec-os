/* io block chaos
add by kfzx-yixy
 */
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/chaosblade-io/chaosblade-exec-os/exec/bin"
	"github.com/chaosblade-io/chaosblade-spec-go/channel"
	"github.com/chaosblade-io/chaosblade-spec-go/util"
	"os"
	"os/exec"
	"path"
	"strings"
)

var UnsupportErr=errors.New("unsupport updateType")
var ReadFileName="/sys/fs/cgroup/blkio/blkio.throttle.read_bps_device"
var WriteFileName="/sys/fs/cgroup/blkio/blkio.throttle.write_bps_device"
var blockIOReadByte, blockIOWriteByte string
var blockIOStart, blockIOStop, blockIONohup bool

func main() {
	flag.StringVar(&blockIOReadByte, "rbyte", "", "read bytes per second")
	flag.StringVar(&blockIOWriteByte, "wbyte", "", "write bytes per second")
	flag.BoolVar(&blockIOStart, "start", false, "start block io")
	flag.BoolVar(&blockIOStop, "stop", false, "stop block io")
	flag.BoolVar(&blockIONohup, "nohup", false, "start by nohup")
	bin.ParseFlagAndInitLog()

	if blockIOStart {
		startBlockIO(blockIOReadByte, blockIOWriteByte)
	} else if blockIOStop {
		stopBlockIO()
	} else if blockIONohup {
		if blockIOReadByte!="" {
			blockRead(blockIOReadByte)
		}
		if blockIOWriteByte!="" {
			blockWrite(blockIOWriteByte)
		}
	} else {
		bin.PrintErrAndExit("less --start or --stop flag")
	}
}

var blockIOBin = "chaos_blockio"
var logFile = util.GetNohupOutput(util.Bin, "chaos_blockio.log")

var cl = channel.NewLocalChannel()

var stopBlockIOFunc = stopBlockIO

// start block io
func startBlockIO(rbytes,wbytes string) {
	ctx := context.Background()
	response := cl.Run(ctx, "nohup",
		fmt.Sprintf(`%s --rbyte=%s --wbyte=%s --nohup=true > %s 2>&1 &`,
			path.Join(util.GetProgramPath(), blockIOBin), rbytes, wbytes, logFile))
	if !response.Success {
		stopBlockIOFunc()
		bin.PrintErrAndExit(response.Err)
		return
	}
	bin.PrintOutputAndExit("success")
}

// echo "8:0 0"> /sys/fs/cgroup/blkio/blkio.throttle.write_bps_device
// echo "8:0 0"> /sys/fs/cgroup/blkio/blkio.throttle.read_bps_device
func stopBlockIO() {
	var errMsg string
	err := updateBlkio("write", "0", context.TODO())
	if err != nil {
		errMsg=fmt.Sprintf("write /sys/fs/cgroup/blkio/blkio.throttle.write_bps_device error: %s ;", err.Error())
	}
	err = updateBlkio("read", "0", context.TODO())
	if err != nil {
		errMsg=fmt.Sprintf("%s write /sys/fs/cgroup/blkio/blkio.throttle.read_bps_device error: %s ;",errMsg, err.Error())
	}
	if errMsg!=""{
		bin.PrintErrAndExit(errMsg)
	}
}

// echo "8:0 1024"> /sys/fs/cgroup/blkio/blkio.throttle.write_bps_device
func blockWrite(bytes string) {
	err := updateBlkio("write", bytes, context.TODO())
	if err != nil {
		bin.PrintAndExitWithErrPrefix(
			fmt.Sprintf("write /sys/fs/cgroup/blkio/blkio.throttle.write_bps_device file error, %s", err.Error()))
	}
}

// echo "8:0 1024"> /sys/fs/cgroup/blkio/blkio.throttle.read_bps_device
func blockRead(bytes string) {
	err := updateBlkio("read", bytes, context.TODO())
	if err != nil {
		bin.PrintAndExitWithErrPrefix(
			fmt.Sprintf("write /sys/fs/cgroup/blkio/blkio.throttle.read_bps_device file error, %s", err.Error()))
	}
}

func getMajMin(ctx context.Context)([]string,error){
	cmd := exec.CommandContext(ctx, "/bin/sh", "-c", `lsblk | awk '$6=="disk"{print $2}'`)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	str := string(output)
	result := strings.Split(str,"\n")
	return result,nil
}

func updateBlkio(updateType,bytes string,ctx context.Context)error{
	var fileName string
	if updateType=="read"{
		fileName =ReadFileName
	}else if updateType=="write"{
		fileName =WriteFileName
	}else{
		return UnsupportErr
	}
	MajMin, err := getMajMin(ctx)
	if err != nil {
		return err
	}
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	for _,mm:=range MajMin {
		if mm=="" {
			continue
		}
		_, err := file.WriteString(fmt.Sprintf("%s %s", mm, bytes))
		if err != nil {
			return err
		}
	}
	return nil
}
