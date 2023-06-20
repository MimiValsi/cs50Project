# CURATOR

#### Video URL: https://youtu.be/Lsn1FYk7ycU

#### Description:
CURATOR is meant to track and centralize defaults found while working.
It's separated in 2 data types:
- Source (place)
- Infos

The program is separated within 4 main folders:

### cmd/
Inside it's the 5 main files which will make the program run.

- handlers send and retrieve data from the http response body/writer. 
    It communicate with database files so the data circulates between PSQL and the browser

- helpers concentrate some web errors to display to the user.
    ex.: 500 or 404

- main file is where the program beggins

- routers file concentrate every page router.

- template file makes sure to create template cache and ensures that
    the files to generate exists

### database/
- errors file has a global error variable to be used when a transaction went wrong

- infos and sources file has every command to insert, update and delete info data

### ui/html/
- base file is the starting point to create a web page

### ui/html/pages/
Each page renders a specific behavior

- home file renders a grid with the source place and a graph to display
    in real time the amount of info per source.

- info..., the 3 files will render their pages to create, update and view
    the info data

- source..., does exactly the same thing as info...

### ui/staic/img
icons and images used are inside

### ui/static/js
This program runs 100% local. the node_modules/ folder regroups JS libs
    used to create and operate this web program

- graph file creates the graph in home page to track the amount of infos
    per source place

- main file create a search in source view page. It search for info status

### ui/static/sass
Every style for the web program are stocked here.

## Tree settings

```
CURATOR/
│
├── cmd/
│   ├── handlers.go
│   ├── helpers.go
│   ├── main.go
│   ├── routers.go
│   └── templates.go
│
├── database/
│   ├── errors.go
│   ├── infos.go
│   └── sources.go
│
├── internal/
│   └── validator/
│       └── validator.go
│
└── ui/
    ├── html/
    │   ├── pages/
    │   │   ├── home.tmpl.html
    │   │   ├── infoCreate.tmpl.html
    │   │   ├── infoUpdate.tmpl.html
    │   │   ├── infoView.tmpl.html
    │   │   ├── sourceCreate.tmpl.html
    │   │   ├── sourceUpdate.tmpl.html
    │   │   └── sourceView.tmpl.html
    │   │
    │   └── base.tmpl.html
    │
    └── static/
        ├── img/
        │   ├── icone_corbeille.png
        │   ├── icone_edition.png
        │   ├── icone_fleche.png
        │   ├── icone_maison.png
        │   ├── icone_ps.png
        │   ├── loupe.png
        │   ├── petit_icone.png
        │   └── searchicon.png
        │
        ├── js/
        │   ├── node_modules/...
        │   ├── graph.js
        │   └── main.js
        └── sass/
            ├── @mdi/...        
            ├── icons/...
            ├── mybulma/...
            └── main.scss
```
