package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
)

var (
	configFilePath = ""
	h = '\u2500'
	last = '\u2516'
	mid = '\u2520'
	v = '\u2503'
	ress = ""

	flags = []cli.Flag{
			&cli.StringFlag{
				Name:        "path",
				Aliases:     []string{"p"},
				Usage: 		"set path to dir, empty value means current directory",
				Destination: &configFilePath,
			},
			&cli.BoolFlag{
				Name:        "no-color",
				Aliases:     []string{"nc"},
				Usage:       "set no color",
			},
			// flag for including dot dirs
			&cli.BoolFlag{
				Name:        "dot-dirs",
				Aliases:     []string{"dd"},
				Usage:       "print include dot directories",
			},
			// flag for including dot files
			&cli.BoolFlag{
				Name:        "dot-files",
				Aliases:     []string{"df"},
				Usage:       "print include dot files",
			},
		}

		counts = Counts{}
)

type Counts struct {
	Dirs 	int64
	Files 	int64
}

var Reset  = "\033[0m"
var Red    = "\033[31m"
var Green  = "\033[32m"
var Yellow = "\033[33m"
var Blue   = "\033[34m"
var Purple = "\033[35m"
var Cyan   = "\033[36m"
var Gray   = "\033[37m"
var White  = "\033[97m"

func init() {
	if runtime.GOOS == "windows" {
		switchColor()
	}
}

func switchColor() {
	Reset  = ""
	Red    = ""
	Green  = ""
	Yellow = ""
	Blue   = ""
	Purple = ""
	Cyan   = ""
	Gray   = ""
	White  = ""
}

func main(){
	app := cli.NewApp()
	app.Commands = cli.Commands{
		&cli.Command{
			Name: "printdir",
			Action: PrintDir,
			Flags: flags,
			Usage: "Show dir tree",
		},
	}
	app.Run(os.Args)
}

func PrintDir(c *cli.Context) error {
	if c.IsSet("no-color") {
		switchColor()
	}

	dotDirs := c.IsSet("dot-dirs")
	dotFiles := c.IsSet("dot-files")

	if configFilePath == "" {
		path, err := os.Getwd()
		if err != nil {
			return err
		}
		configFilePath = path
	}
	print(configFilePath,0, dotDirs, dotFiles)
	lis := strings.Split(ress,"\n")
	lis2 := make([][]string,0)
	for _,c := range lis {
		lis2 = append(lis2,strings.Split(c,""))
	}
	for i:=0;i<len(lis2);i++  {
		for j:=0;j<len(lis2[i]);j++{
			if i<len(lis2)-2 && (lis2[i][j] ==string(mid) || lis2[i][j]==string(v)) && lis2[i+1][j]==" " {
				lis2[i+1][j]=string(v)
			}
		}
	}
	ress = ""
	for _,li := range lis2 {
		ress+= strings.Join(li,"")+"\n"
	}
	fmt.Print(ress)
	fmt.Println(fmt.Sprintf("%vDirectories: %d\nFiles: %d%v", Green, counts.Dirs, counts.Files, Reset))
	return nil
}

func print(s string, n int, dotDirs, dotFiles bool) {
	dir,er := ioutil.ReadDir(s)
	if er == nil {
		for i,d := range dir {
			res := ""
			for j:=0;j<n;j++{
				res+=" "
			}
			if i==len(dir)-1 {
				res += string(last)
			}else{
				res += string(mid)
			}
			if d.IsDir() {
				if !dotDirs && d.Name()[0] == 46 {
					continue
				}
				res+= Blue + string(h)+d.Name() + Reset
				ress+=res+"\n"
				counts.Dirs += 1
			}else{
				if !dotFiles && d.Name()[0] == 46 {
					continue
				}
				res+= Yellow + string(h)+string(h)+d.Name() + Reset
				ress+=res+"\n"
				counts.Files += 1
			}
			print(s+"/"+d.Name(),n+1, dotDirs, dotFiles)
		}
	}
}



