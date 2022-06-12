package files

import (
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

func GetFiles(ctx *gin.Context) {
	imgs, err := GetImages(ctx)
	if err != nil {
		ctx.JSON(500, err)
		return
	}
	ctx.JSON(200, imgs)
}

func GetThumbnail(ctx *gin.Context) {
	id := ctx.Param("hash")
	img, err := GetImageFile(ctx, id)
	if err != nil {
		ctx.JSON(500, err)
		return
	}
	imgb, err := img.Thumbnail()
	if err != nil {
		ctx.JSON(500, err)
		return
	}
	buf, err := ioutil.ReadAll(imgb)
	if err != nil {
		ctx.JSON(500, err)
		return
	}
	ctx.Data(200, "image/"+img.Type, buf)
}

func GetImage(ctx *gin.Context) {
	id := ctx.Param("hash")
	img, err := GetImageFile(ctx, id)
	if err != nil {
		ctx.JSON(500, err)
		return
	}
	imgb, err := img.Image()
	if err != nil {
		ctx.JSON(500, err)
		return
	}
	buf, err := ioutil.ReadAll(imgb)
	if err != nil {
		ctx.JSON(500, err)
		return
	}
	ctx.Header("Cache-Control", "max-age=31536000")
	ctx.Data(200, "image/"+img.Type, buf)
}

func UploadFiles(ctx *gin.Context) {
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(500, err)
		return
	}
	fmt.Println(form.File)
	newFiles := make([]*file, 0)
	for _, f := range form.File["files"] {
		fd, err := f.Open()
		if err != nil {
			ctx.JSON(500, err)
			return
		}
		f_obj, err := NewImage(ctx.Request.Context(), f.Filename, fd)
		if err != nil {
			ctx.JSON(500, err)
			return
		}
		err = f_obj.Save(ctx)
		if err != nil {
			ctx.JSON(500, err)
			return
		}
		newFiles = append(newFiles, f_obj)
	}
	ctx.JSON(200, newFiles)
}
