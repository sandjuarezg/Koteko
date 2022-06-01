const elements = document.querySelectorAll(".block--products")

// Scroll
window.onscroll = function() {
    var y = window.scrollY;

    elements.forEach(function(element) {
        var elementTop = element.offsetTop - 400
        
        if (y > elementTop) {
            element.style.left = "0";
        }
    })
}

function showPassword(){
    var tipo = document.getElementById("password");
    if(tipo.type == "password"){
        tipo.type = "text";
    }else{
        tipo.type = "password";
    }
}