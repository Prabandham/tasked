$(document).on("turbolinks:load", function() {
    // Hiding Links
    var path = window.location.pathname;
    
    if(path == "/") {
        $("#login-link").removeClass("hidden")
    } else {
        $("#login-link").addClass("hidden")
    }

    // Task completion for List View
    $(".task-check").on("click", function() {
        var task = $(this);
        $.post("/complete_task", { value: task.attr("id").split("task_")[1] })
            .done(function(data) {
                task.parent().removeClass("list-group-item-warning");
                task.parent().addClass(data.cssClass);
                task.remove();
            });
    });

    //Remove alert after 3 seconds.
    $(".alert").delay(3000).fadeOut("slow");

    // Add user to Project
    $("#add_user").on("click", function(e) {
        $("#user_email").toggleClass("hidden");
        e.preventDefault();
    });

    // Enter on adding User.
    $("#user_email").keypress(function(e) {
        if(e.which == 13) {
            email = $(this).val();
            console.log(email);
            $.post("/add_user", { project_id: path.split("/")[2], email: email })
                .done(function(data) {
                    if(data.status == "success") {
                        console.log("Came here");
                        window.location.reload()
                    }
                    else if(data.status == "failure") {
                        console.log("Came to failure");
                        alert(data.message);
                    }
                })
                .error(function() {
                    alert("Some thing went wrong !! Please Try again !!");
                });
        }
    });

    $("#NewPhaseName").on("click", function() {
        $(this).html("");
    });

    // Foucus out event for all contenteditable fields
    $('[contenteditable]').on('focusout', function(e){
        var sr = window.getSelection().getRangeAt(0).commonAncestorContainer;
        var el = sr.parentNode;
        if($(el).attr("id") != null || $(el).attr("id") != undefined) {
            var el_type = $(el).attr("id").split("_")[0];
            var el_id = $(el).attr("id").split("_")[1];
            var el_value = $(el).context.innerHTML;
        } else {
            $("#NewPhaseName").html("Add a New Phase ..");
            return
        }
        if(el_type == "taskName") {
            $.post("/update_task", { task_id: el_id, task_name: el_value })
        }
        if(el_type == "phaseName") {
            $.post("/update_phase", { phase_id: el_id, phase_name: el_value })
        }
        if(el_type == "NewPhaseName")  {
            if(el_value != "") {
                $.post("/add_phase", { project_id: path.split("/")[2], phase_name: el_value })
                    .success(function(data) {
                        if(data.status == "success") {
                            window.location.reload()
                        }
                    })
            } else {
                $("#NewPhaseName").html("Add a New Phase ..");
            }
        }
    });
});
