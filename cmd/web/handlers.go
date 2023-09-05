package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"Snippetbox.Shreshth/internal/models"
	"Snippetbox.Shreshth/internal/validator"
	"github.com/julienschmidt/httprouter"
)


type snippetCreateForm struct {
	Title string `form:"title"`
	Content string `form:"content"`
	Expires int `form:"expires"`
	validator.Validator `form:"-"`
	
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	
	
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets
	// Use the new render helper.
	app.render(w, http.StatusOK, "home.tmpl", data)


}
	func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	
		params := httprouter.ParamsFromContext(r.Context())
		// We can then use the ByName() method to get the value of the "id" named
		// parameter from the slice and validate it as normal.
		id, err := strconv.Atoi(params.ByName("id"))
		if err != nil || id < 1 {
		app.notFound(w)
		return
		}
		snippet, err := app.snippets.Get(id)
		if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
		app.notFound(w)
		} else {
		app.serverError(w, err)
		}
		return
		}
		flash := app.sessionManager.PopString(r.Context(), "flash")
		data := app.newTemplateData(r)
		data.Snippet = snippet
		data.Flash = flash
		app.render(w, http.StatusOK, "view.tmpl", data)
}






func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	// Initialize a new createSnippetForm instance and pass it to the template.
	// Notice how this is also a great opportunity to set any default or
	// 'initial' values for the form --- here we set the initial value for the
	// snippet expiry to 365 days.
		data.Form = snippetCreateForm{
		Expires: 365,
	}
	app.render(w, http.StatusOK, "create.tmpl", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	var form snippetCreateForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// Declare a new empty instance of the snippetCreateForm struct.
	
	// Call the Decode() method of the form decoder, passing in the current
	// request and *a pointer* to our snippetCreateForm struct. This will
	// essentially fill our struct with the relevant values from the HTML form.
	// If there is a problem, we return a 400 Bad Request response to the client.
	err = app.formDecoder.Decode(&form, r.PostForm)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// Then validate and use the data as normal...
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}
	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)

}