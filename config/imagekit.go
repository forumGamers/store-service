package config

import (
	"errors"
	"os"

	"context"

	"github.com/codedius/imagekit-go"
	"github.com/joho/godotenv"
)

func ImagekitConnection() *imagekit.Client{
	if err := godotenv.Load() ; err != nil {
		panic(err.Error())
	}

	IMAGEKIT_PRIVATE_KEY := os.Getenv("IMAGEKIT_PRIVATE_KEY")
	IMAGEKIT_PUBLIC_KEY := os.Getenv("IMAGEKIT_PUBLIC_KEY")

	opts := imagekit.Options{
		PublicKey: IMAGEKIT_PUBLIC_KEY,
		PrivateKey: IMAGEKIT_PRIVATE_KEY,
	}

	if ik,err := imagekit.NewClient(&opts) ; err != nil {
		panic(err.Error())
	}else {
		return ik
	}
}

func UploadImage(file []byte,fileName string,folder string) (string, string , error){

	ur := imagekit.UploadRequest{
		File: file,
		FileName: fileName,
		UseUniqueFileName: true,
		Folder: folder,
	}

	ctx := context.Background()

	if upr,err := ImagekitConnection().Upload.ServerUpload(ctx,&ur) ; err != nil {
		return "",
			   "",
			   errors.New(err.Error())
	}else {
		return upr.URL,
			   upr.FileID,
			   nil
	}
}

func UpdateImage(file []byte,fileName string,folder string,updatedfileId string) (string, string, error) {

	if url,fileId,err := UploadImage(file,fileName,folder) ; err != nil {
		return "","",errors.New(err.Error())
	}else {
		if updatedfileId != "" {
			ctx := context.Background()

			if err := ImagekitConnection().Media.DeleteFile(ctx,updatedfileId) ; err != nil {
				return url,fileId,err
			}
		}
		return url,fileId,nil
	}
}