var listsDB = [];
var todosDB = [];

//-> Icons

//Check circle icon
function checkCircle(){
    return feather.icons["check-circle"].toSvg({ width: 24, height: 24, class: "me-2" });
}

//-> Script to remember the scroll position and download todos and lists from the server
document.addEventListener("DOMContentLoaded", function (event) {
    download();
    document.getElementById("search-keyword").value = "";
    var scrollpos = localStorage.getItem('scrollpos');
    if (scrollpos) window.scrollTo(0, scrollpos);
});

window.onbeforeunload = function (e) {
    localStorage.setItem('scrollpos', window.scrollY);
};

//-> Api calls

// Function to download lists and todos from the server and update the page
function download(){
    fetch("api/lists/get", {
        method: "GET",
        headers: {'Content-Type': 'application/json'},
    }).then(res => {
        if(res.status == 200){
            res.json().then(out =>{
                if(out != null && out != undefined && out != ""){
                    listsDB = out;
                }
            })
        } else {
            res.text().then(text => {
                console.error("Errore " + res.status + " - " + text);
                showNotification("Errore " + res.status + " - " + text, "error");
            });
        }
    });

    fetch("api/todos/get", {
        method: "GET",
        headers: {'Content-Type': 'application/json'},
    }).then(res => {
        if(res.status == 200){
            res.json().then(out =>{
                if(out != null && out != undefined && out != ""){
                    //Save only not completed todos in the local "database"
                    out.forEach(todo => {
                        if (todo.isCompleted == false) {
                            todosDB.push(todo);
                        }
                    });

                    //Order by modified date
                    todosDB.sort(function(a, b){
                        return new Date(b.lastModified) - new Date(a.lastModified);
                    });
                }
                search(); //Update the page with the new data
                showNotification(checkCircle() + "Dati ottenuti con successo!");
            })
        } else {
            res.text().then(text => {
                console.error("Errore " + res.status + " - " + text);
                showNotification("Errore " + res.status + " - " + text, "error");
            });
        }
    });
}

//-> Search function
//Search on button click
document.getElementById("search-button").addEventListener("click", function(event) {
    event.preventDefault();
    search();
});

// Real-time search during typing
document.getElementById("search-keyword").addEventListener("keyup", function(event) {
    event.preventDefault();
    search();
});

// Function to empty list
function emptyList(){
    document.getElementById("list").textContent = '';
}

// Local search function (also used to display all elements)
function search() {
    //Convert input to lowercase as the search is case insensitive
    let input = document.getElementById("search-keyword").value.toLowerCase();
    let searchField = parseInt(document.getElementById("search-field").value);
    let searchQuery = [];
    let list = []; //List of elements to put in the html page
    
    //Check if the input is empty
    if(input != "" && input != null && input != undefined){
        switch (searchField){
            case 1: //Filter by title
                searchQuery = todosDB.filter(todo => todo.title.toLowerCase().includes(input));
                
                /* //Sort by first occurrence of the search string? Now it is sorted by last modified date
                searchQuery.sort(function(a, b){
                    return a.title.toLowerCase().indexOf(input) - b.title.toLowerCase().indexOf(input);
                });*/

            case 2: //Filter by description
                searchQuery = todosDB.filter(todo => todo.description.toLowerCase().includes(input));

            case 3: //Filter by title or description
                searchQuery = todosDB.filter(todo => todo.title.toLowerCase().includes(input) || todo.description.toLowerCase().includes(input));
        }

        //Common code for all filters
        if (searchQuery.length > 0){
            searchQuery.forEach(todo => {
                list.push(createListItem(todo));
            });
            document.getElementById("list").innerHTML = list.join('');
        } else {
            document.getElementById("list").textContent = 'Nessun elemento corrisponde al criterio di ricerca';
        }

    } else { //No input, show all elements
        if (todosDB.length > 0){
            todosDB.forEach(todo => {
                list.push(createListItem(todo));
            });
            document.getElementById("list").innerHTML = list.join('');
        } else {
            document.getElementById("list").textContent = 'Nessun elemento presente nel database';
        }
    }
}

//createListItem creates the html of a list item for the search function
function createListItem(todo){
    let style = "";
    let formattedList = '<span class="badge text-bg-dark bg-opacity-75 fw-normal rounded-3 me-auto mt-1">Generale</span>';
    if (todo.listId != "") {
        let list = listsDB.find(x => x.id == todo.listId);
        style = ' style="background-color: rgb(' + list.color[0] + "," + list.color[1] + "," + list.color[2] + ') !important"';
        formattedList = '<span class="badge text-bg-dark bg-opacity-75 fw-normal rounded-3 me-auto mt-1"' + style + ">" + list.name + "</span>";
    }
    
    return `<a id="${todo.id}" class="list-group-item list-group-item-action">
    <div class="d-flex flex-row justify-content-between">
        <div id="${todo.id}-left" class="flex-grow-4 flex-grow-1 d-flex flex-row">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="1.52rem" height="1.52rem" min-width="1.52rem" min-height="1.52rem" class="main-grid-item-icon pe-1 flex-shrink-0" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2">
                <circle cx="12" cy="12" r="10" />
            </svg>
            <div class="d-flex flex-column">
                <div class="d-flex flex-column">
                    <h5 id="${todo.id}-title" class="d-flex align-items-top mb-auto fs-5">${todo.title}</h5>
                    <p id="${todo.id}-description" class="m-0 fs-6">${todo.description}</p>
                </div>
                ${formattedList}
            </div>
        </div>
    </div>    
</a>`
}

//-> Notification toast
const notifyToast = new bootstrap.Toast('#notifyToast', {});

/*
    showNotification shows a notification toast with the specified text and type on top of the page.

    Input fields:
    - text (string) -> Notification text
    - type (string) -> Notification type (success, warning, error), default: success
*/
function showNotification(text, type){
    if (typeof(text) != "string" || text == "" || text == null || text == undefined){
        console.error("Errore: testo non valido. text: " + text + " - Tipo di 'text': " + typeof(text));
        return;
    }
    switch (type){
        case "warning":
            document.getElementById("notifyToast").classList.remove("text-bg-success", "text-bg-danger", "text-bg-warning");
            document.getElementById("notifyToast").classList.add("text-bg-warning");
            break;
        case "error":
            document.getElementById("notifyToast").classList.remove("text-bg-success", "text-bg-danger", "text-bg-warning");
            document.getElementById("notifyToast").classList.add("text-bg-danger");
            break;
        default:
            document.getElementById("notifyToast").classList.remove("text-bg-success", "text-bg-danger", "text-bg-warning");
            document.getElementById("notifyToast").classList.add("text-bg-success");
            break;
    }
    document.getElementById("notifyToastBody").innerHTML = text;
    notifyToast.show();
}