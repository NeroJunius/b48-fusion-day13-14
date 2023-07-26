package main

import (
	"batch48/connection"
	"context"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type Projects struct {
	ID          int
	ProjectName string
	Author      string

	StartDateFormat string
	EndDateFormat   string
	DurationFormat  string

	DescriptionProject string
	NodeJS             bool
	ReactJS            bool
	NextJS             bool
	TypeScript         bool
	Img                string

	StartDate time.Time
	EndDate   time.Time
	Duration  time.Duration
}


func main() {
	connection.DatabaseConnect()
	e := echo.New()
	e.Static("/assets", "assets")

	e.GET("/", Home)
	e.GET("/contactMe", contactMe)
	e.GET("/testimonial", testimonials)
	e.GET("/createProject", createProject)
	e.GET("/projectDetail/:id", projectDetail)
	e.GET("/editProject/:id", editProject)

	e.POST("/add-project", AddProject)
	e.POST("/edit-project/:id", EditedProject)
	e.POST("/delete-project/:id", DeleteProject)

	fmt.Println("server started on port 5900")
	e.Logger.Fatal(e.Start("localhost:5900"))
}

// List Fungsi GET Project nya /

func Home(c echo.Context) error {
	data, _ := connection.Conn.Query(context.Background(), "SELECT id, project_title, start_date, end_date, description, node_js, react_js, next_js, type_script, image FROM tb_projects;")
	var results []Projects
	for data.Next() {
		each := Projects{}

		err := data.Scan(&each.ID, &each.ProjectName, &each.StartDate, &each.EndDate, &each.DescriptionProject, &each.NodeJS, &each.ReactJS, &each.NextJS, &each.TypeScript, &each.Img)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}

		each.Duration = each.EndDate.Sub(each.StartDate)
		each.DurationFormat = DurationFormat(each.Duration)
		results = append(results, each)
	}

	tmpl, err := template.ParseFiles("tabs/index.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	projects := map[string]interface{}{
		"Project": results,
	}

	return tmpl.Execute(c.Response(), projects)
}

func contactMe(c echo.Context) error {
	var tmpl, err = template.ParseFiles("tabs/contact.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	return tmpl.Execute(c.Response(), nil)
}

func createProject(c echo.Context) error {
	var tmpl, err = template.ParseFiles("tabs/project.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	return tmpl.Execute(c.Response(), nil)
}

func testimonials(c echo.Context) error {
	var tmpl, err = template.ParseFiles("tabs/testimonial.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	return tmpl.Execute(c.Response(), nil)
}

func projectDetail(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	var ListProjects = Projects{}

	err := connection.Conn.QueryRow(context.Background(),
		" SELECT id, project_title, start_date, end_date, description, node_js, react_js, next_js, type_script, image FROM tb_projects WHERE id=$1", id).Scan(
		&ListProjects.ID, &ListProjects.ProjectName, &ListProjects.StartDate, &ListProjects.EndDate, &ListProjects.Duration, &ListProjects.DescriptionProject,
		&ListProjects.NodeJS, &ListProjects.ReactJS, &ListProjects.NextJS, &ListProjects.TypeScript, &ListProjects.Img)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	
	data := map[string]interface{}{
		"Project":   ListProjects,
		"StartDate": ListProjects.StartDate.Format("June 28, 1999"),
		"EndDate":   ListProjects.EndDate.Format("August 28, 2023"),
	}

	var tmpl, errTemp = template.ParseFiles("tabs/project-detail.html")
	if errTemp != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), data)
}

func editProject(c echo.Context) error {

	id, _ := strconv.Atoi(c.Param("id"))

	var Previous_Data = Projects{}

	err := connection.Conn.QueryRow(context.Background(),"SELECT id, project_title, start_date, end_date, description, node_js, react_js, next_js, type_script, image FROM tb_projects WHERE id=$1", id ).Scan(&Previous_Data.ID, &Previous_Data.ProjectName, &Previous_Data.StartDate, &Previous_Data.EndDate, &Previous_Data.DescriptionProject, &Previous_Data.NodeJS, &Previous_Data.ReactJS, &Previous_Data.NextJS, &Previous_Data.TypeScript, &Previous_Data.Img )

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	data := map[string]interface{}{
		"Previous_Data": Previous_Data,
	}

	tmpl, errTemp := template.ParseFiles("tabs/edit-form.html")

	if errTemp != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	return tmpl.Execute(c.Response(), data)
}

// LIST QUERY PROJECT //

// time //
func DurationFormat(Duration time.Duration) string {
	if Duration <= 24*time.Hour {
		return "Less than a day"
	}

	Days := int(Duration.Hours() / 24)
	Weeks := Days / 7
	Months := Days / 30
	Years := Months / 12

	if Years > 1 {
		return fmt.Sprintf("%d years", Years)
	} else if Years == 1 {
		return "A year"
	} else if Months > 1 {
		return fmt.Sprintf("%d months", Months)
	} else if Months == 1 {
		return "A month"
	} else if Weeks > 1 {
		return fmt.Sprintf("%d weeks", Weeks)
	} else if Weeks == 1 {
		return "A week"
	} else if Days > 1 {
		return fmt.Sprintf("%d days", Days)
	} else {
		return "A day"
	}
}

// buat project nya //
func AddProject(c echo.Context) error {
	ProjectName := c.FormValue("projectName")
	StartDate := c.FormValue("startDate")
	EndDate := c.FormValue("endDate")
	DescriptionProject := c.FormValue("projectDescription")

	var NodeJS bool
	if c.FormValue("nodeJS") == "yes" {
		NodeJS = true
	}
	var NextJS bool
	if c.FormValue("nextJS") == "yes" {
		NextJS = true
	}
	var ReactJS bool
	if c.FormValue("reactJS") == "yes" {
		ReactJS = true
	}
	var TypeScript bool
	if c.FormValue("typeScript") == "yes" {
		TypeScript = true
	}

	Img := c.FormValue("imageProject")

	_, err := connection.Conn.Exec(context.Background(),
		"INSERT INTO tb_projects (project_title, start_date, end_date, description, node_js, react_js, next_js, type_script, image) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
		ProjectName, StartDate, EndDate, DescriptionProject, NodeJS, ReactJS, NextJS, TypeScript, Img)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	fmt.Println(ProjectName, StartDate, EndDate, DescriptionProject, NodeJS, ReactJS, NextJS, TypeScript)


	// dataProjects = append(dataProjects, createProject)

	return c.Redirect(http.StatusMovedPermanently, "/")
}

// edited project //
func EditedProject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	ProjectName := c.FormValue("projectName")
	StartDate := c.FormValue("startDate")
	EndDate := c.FormValue("endDate")
	DescriptionProjects := c.FormValue("projectDescription")

	var NodeJS bool
	if c.FormValue("nodeJS") == "yes" {
		NodeJS = true
	}
	var NextJS bool
	if c.FormValue("nextJS") == "yes" {
		NextJS = true
	}
	var ReactJS bool
	if c.FormValue("reactJS") == "yes" {
		ReactJS = true
	}
	var TypeScript bool
	if c.FormValue("typeScript") == "yes" {
		TypeScript = true
	}


	_, err := connection.Conn.Exec(context.Background(),
		"UPDATE tb_projects SET project_title=$1, start_date=$2, end_date=$3, description=$4, node_js=$5, react_js=$6, next_js=$7, type_script=$8, image=$9) WHERE id=$1",
		ProjectName, StartDate, EndDate, DescriptionProjects, NodeJS, ReactJS, NextJS, TypeScript, "Kal S3.png", id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	// redirectURL := fmt.Sprintf("/?id=#project", id)
	return c.Redirect(http.StatusMovedPermanently, "/")
}

// delete project //
func DeleteProject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	// dataProjects = append(dataProjects[:id], dataProjects[id+1:]...)

	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM tb_projects WHERE id=$1", id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	fmt.Println("Berhasil menghapus project!")

	return c.Redirect(http.StatusMovedPermanently, "/")
}
