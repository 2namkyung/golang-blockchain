(function($){
    'use strict';
    $(function(){

        var MakeQR = function(qrbytes){
            document.getElementById("qrcode").src = "data:image/png;base64, " + qrbytes;
        }
        
        $('.qrcode-btn').on("click", function(event){
            event.preventDefault();

            $.get("/qrcode", function(qrbytes){
                MakeQR(qrbytes);
            })
        });

        
    });
})(jQuery);