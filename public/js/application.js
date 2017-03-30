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
});
