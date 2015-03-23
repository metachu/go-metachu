function errorToast(title,detail) {
	toast('<div><h5 class="red-text">'+title+": "+detail+'</h5></div>', 4000) // 4000 is the duration of the toast
}

function successToast(title,detail,refresh) {
	if (refresh) {
		toast('<div><h5 class="green-text">'+title+": "+detail+'</h5></div>', 1000,'',function(){location.reload(true)}); // 4000 is the duration of the toast
	} else {
		toast('<div><h5 class="green-text">'+title+": "+detail+'</h5></div>', 4000) // 4000 is the duration of the toast
	}
}

$(document).ready(function(){
	var fileRows = $('.file-row');
	var currentFile = null;
	
	// Setting up modal form actions
	$("#zipModalZipAndDelete")[0].addEventListener('click', function(){
		jQuery.post("/files/json/zip/",data={"action":"zipAndDelete","filepath":currentFile.attr("data-rel-path"), "newname":$('#zipFileName')[0].value})
		.done(function(data) {
			successToast(data["data"]["title"],data["data"]["detail"],true);

		})
		.fail(function(data){
			errorToast(data["responseJSON"]["errors"]["title"],data["responseJSON"]["errors"]["detail"]);
		});

	});
	$("#zipModalZip")[0].addEventListener('click', function(){
		jQuery.post("/files/json/zip/",data={"action":"zip","filepath":currentFile.attr("data-rel-path"), "newname":$('#zipFileName')[0].value})
		.done(function(data) {
			successToast(data["data"]["title"],data["data"]["detail"],true);
		})
		.fail(function(data){
			errorToast(data["responseJSON"]["errors"]["title"],data["responseJSON"]["errors"]["detail"]);
		});
	});


	$("#delModalConfirm")[0].addEventListener('click', function(){
		jQuery.post("/files/json/delete/",data={"filepath":currentFile.attr("data-rel-path")})
		.done(function (data){
			currentFile.hide();
			successToast(data["data"]["title"],data["data"]["detail"]);
		})
		.fail(function (data){
			errorToast(data["responseJSON"]["errors"]["title"],data["responseJSON"]["errors"]["detail"]);
		});
	});

	$("#movModalConfirm")[0].addEventListener('click', function(){

	});
	
	// Setting up modal display buttons for each row
	fileRows.each(function() {
		var thisRow = $(this);
		var fileName = $(this).attr("data-filename");
		var md5 = $(this).attr("data-md5");
		$(this).find(".zip-button").each(function() {
			this.addEventListener('click', function() {
				var modal = $('#zipModal');
				var zipFileInput = $("#zipFileName");
				currentFile = thisRow;
				$("#zipModalForm")[0].reset();
				zipFileInput.attr("value",fileName+".zip");
				modal.openModal();
			}, false);
		});

		$(this).find(".del-button").each(function() {
			this.addEventListener('click', function() {
				currentFile = thisRow;
				var modal = $('#delModal');
				modal.openModal();
			}, false);
		});

		$(this).find(".mov-button").each(function() {
			this.addEventListener('click', function() {
				currentFile = thisRow;
				var modal = $('#movModal');
				modal.openModal();
			}, false);
		});		

	});
    // $('#modal1').openModal();
	$('.skew-title').children('span').hover(function() {
		
		var n = $(this).index();
		n++;
		$(this).addClass('flat');
		
		if ((n % 2 == 0)) {  
			$(this).prev().addClass('flat');
		} else {
			if (!$(this).hasClass('last')) {
				$(this).next().addClass('flat');
			}  
		}
		
	}, function() {
		
		$('.flat').removeClass('flat');
		
	});
});



// myEl.addEventListener('click', function() {
//     alert('Hello world');
// }, false);



