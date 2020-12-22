package main 
//env GOOS=linux go build
import (
	"fmt"
	"io/ioutil"
	"strconv"
	"log"
	"archive/zip"
	"io"
	"net/http"
	"os"
	"encoding/json"
	"time"
	"sync"
    "runtime"
)

type dataStruct struct {
	ArchiveName string `json:"archive_name"`
	OutputPath string `json:"output_path"`
	Data folderStruct `json:"data"`
}

type folderSourceStruct struct {
	Name string `json:"name"`
	URL string `json:"url"`
	Extention string `json:"extention"`
}

type folderStruct struct {
	Name string `json:"name"`
	Sources []folderSourceStruct `json:"sources"`
	Folders []folderStruct `json:"folders"`
}

type sourceStruct struct {
	Name string
	Extention string
	URL string
	Folder string
}


func init(){

}

func main() {
    runtime.GOMAXPROCS(4)

	var wg sync.WaitGroup
	var oSource []sourceStruct

	sTempFolder := "./temp/" + makeTimestamp()

	sJSON := os.Args[1]
	
	oData := dataStruct{}
	json.Unmarshal([]byte(sJSON), &oData)

	startConstructStructure(&wg,&oSource,sTempFolder,oData.Data)
	downloadSources(&wg,&oSource)
	wg.Wait()

	createFolder(oData.OutputPath)
	zipWriter(sTempFolder,oData.OutputPath+"/"+oData.ArchiveName+".zip")
	os.RemoveAll(sTempFolder)

	print(true)
}

func startConstructStructure(wg *sync.WaitGroup, oSource *[]sourceStruct, sFolder string, oData folderStruct) {
	sFolderNew := sFolder+"/"+oData.Name
	createFolder(sFolderNew)

	for _,oFile := range oData.Sources {
			*oSource = append(*oSource,sourceStruct{Name:oFile.Name,URL:oFile.URL,Folder:sFolderNew+"/",Extention:oFile.Extention})
	}

	if oData.Folders != nil {
		constructStructure(wg,oSource,sFolderNew,oData.Folders)
	}
}

func constructStructure(wg *sync.WaitGroup, oSource *[]sourceStruct, sFolder string, oData []folderStruct) {

	for _,oFolder := range oData {
		sFolderNew := sFolder+"/"+oFolder.Name
		createFolder(sFolderNew)

		for _,oFile := range oFolder.Sources {
			*oSource = append(*oSource,sourceStruct{Name:oFile.Name,URL:oFile.URL,Folder:sFolderNew+"/",Extention:oFile.Extention})

		}

		if oFolder.Folders != nil {
			constructStructure(wg,oSource,sFolderNew,oFolder.Folders)
		}
	}
}

func downloadSources(wg *sync.WaitGroup,oData *[]sourceStruct){

	for _,oSource := range *oData {
		wg.Add(1)
		go downloadFile(wg,oSource.Folder+oSource.Name+"."+oSource.Extention,oSource.URL)
	}
}


func downloadFile(wg *sync.WaitGroup,filepath string, url string) error {
	defer wg.Done()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err

}

func createFolder(sPath string){
	_, err := os.Stat(sPath)
 
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(sPath, 0777)
		if errDir != nil {
			log.Fatal(err)
		}
 
	}
}

func zipWriter(baseFolder, sOutput string) {

    outFile, err := os.Create(sOutput)
    if err != nil {
        fmt.Println(err)
    }
    defer outFile.Close()

    w := zip.NewWriter(outFile)

    addFiles(w, baseFolder+"/", "")

    if err != nil {
        fmt.Println(err)
    }

    err = w.Close()
    if err != nil {
        fmt.Println(err)
    }
}

func addFiles(w *zip.Writer, basePath, baseInZip string) {

    files, err := ioutil.ReadDir(basePath)
    if err != nil {
        fmt.Println(err)
    }

    for _, file := range files {
        if !file.IsDir() {
            dat, err := ioutil.ReadFile(basePath + file.Name())
            if err != nil {
                fmt.Println(err)
            }

            f, err := w.Create(baseInZip + file.Name())
            if err != nil {
                fmt.Println(err)
            }
            _, err = f.Write(dat)
            if err != nil {
                fmt.Println(err)
            }
        } else if file.IsDir() {

            newBase := basePath + file.Name() + "/"

            addFiles(w, newBase, baseInZip  + file.Name() + "/")
        }
    }
}

func makeTimestamp() string {
    return strconv.FormatInt(time.Now().UnixNano() / int64(time.Millisecond),10)
}