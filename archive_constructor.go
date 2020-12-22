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

	StartConstructStructure(&wg,&oSource,sTempFolder,oData.Data)
	DownloadSources(&wg,&oSource)
	wg.Wait()

	CreateFolder(oData.OutputPath)
	ZipWriter(sTempFolder,oData.OutputPath+"/"+oData.ArchiveName+".zip")
	os.RemoveAll(sTempFolder)

	print(true)
}

func StartConstructStructure(wg *sync.WaitGroup, oSource *[]sourceStruct, sFolder string, oData folderStruct) {
	sFolderNew := sFolder+"/"+oData.Name
	CreateFolder(sFolderNew)

	for _,oFile := range oData.Sources {
			*oSource = append(*oSource,sourceStruct{Name:oFile.Name,URL:oFile.URL,Folder:sFolderNew+"/",Extention:oFile.Extention})
	}

	if oData.Folders != nil {
		ConstructStructure(wg,oSource,sFolderNew,oData.Folders)
	}
}

func ConstructStructure(wg *sync.WaitGroup, oSource *[]sourceStruct, sFolder string, oData []folderStruct) {

	for _,oFolder := range oData {
		sFolderNew := sFolder+"/"+oFolder.Name
		CreateFolder(sFolderNew)

		for _,oFile := range oFolder.Sources {
			*oSource = append(*oSource,sourceStruct{Name:oFile.Name,URL:oFile.URL,Folder:sFolderNew+"/",Extention:oFile.Extention})

		}

		if oFolder.Folders != nil {
			ConstructStructure(wg,oSource,sFolderNew,oFolder.Folders)
		}
	}
}

func DownloadSources(wg *sync.WaitGroup,oData *[]sourceStruct){

	for _,oSource := range *oData {
		wg.Add(1)
		go DownloadFile(wg,oSource.Folder+oSource.Name+"."+oSource.Extention,oSource.URL)
	}
}


func DownloadFile(wg *sync.WaitGroup,filepath string, url string) error {
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

func CreateFolder(sPath string){
	_, err := os.Stat(sPath)
 
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(sPath, 0777)
		if errDir != nil {
			log.Fatal(err)
		}
 
	}
}

func ZipWriter(baseFolder, sOutput string) {

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