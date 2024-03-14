package controllers

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"log"
	"log/slog"
	"springoff/internal/models/album"
)

type Album struct {
	album *album.Album
}

//type UploadAlbum struct {
//	Title string `form:"title"`
//	Cover string `form:"cover"`
//	Album string `json:"album"`
//}

func AlbumNew(db *sql.DB) *Album {
	alb := album.New(db)
	return &Album{album: alb}
}

func (a *Album) GetAlbum() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		idAlbum := c.Params("idAlbum")
		albumImages, err := a.album.GetAlbumImages(idAlbum)
		if err != nil {
			slog.Error("Error get album images", "error", err, "ID album", idAlbum)
			return c.SendStatus(404)
		}

		title, err := a.album.GetTitleAlbum(idAlbum)
		if err != err {
			slog.Error("Error get title album", "error", err, "id album", idAlbum)
		}

		return c.Render("album",
			fiber.Map{"AlbumImages": albumImages,
				"Title": title})
	}
}

func (a *Album) UploadAlbum() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {

		//upload := new(UploadAlbum)
		//
		//if err := c.BodyParser(&upload); err != nil {
		//	log.Print(upload)
		//}
		form, err := c.MultipartForm()
		// Parse the multipart form:
		if err == nil {
			// => *multipart.Form
			//title := form.Value["title"][0]
			//Validate title
			//if len(title) > 0 {
			//........
			//}
			title := form.Value["title"]
			cover := form.File["cover"]
			albumImage := form.File["albumUpload[]"]

			log.Print(cover, albumImage, title)
			// Loop through files:
			//for _, file := range files {
			//	fmt.Println(file.Filename, file.Size, file.Header["Content-Type"][0])
			//	// => "tutorial.pdf" 360641 "application/pdf"
			//
			//	// Save the files to disk:
			//	if err := c.SaveFile(file, fmt.Sprintf("./%s", file.Filename)); err != nil {
			//		return err
			//	}
			//}
		}

		return nil
	}
}
