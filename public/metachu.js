
$(document).ready(function(){
	var fileRows = $('.file-row');
	
	// Setting up modal form actions
	$("#zipModalCancel")[0].addEventListener('click', function(){
		jQuery.post("/files/json/zip/",data="{test:'lol'}");
	});

	$("#zipModalZipAndDelete")[0].addEventListener('click', function(){

	});
	$("#zipModalZip")[0].addEventListener('click', function(){

	});
	$("#delModalCancel")[0].addEventListener('click', function(){

	});
	$("#delModalConfirm")[0].addEventListener('click', function(){

	});
	$("#movModalCancel")[0].addEventListener('click', function(){

	});
	$("#movModalConfirm")[0].addEventListener('click', function(){

	});
	
	
	


	// Setting up modal display buttons for each row
	fileRows.each(function() {
		var fileName = $(this).attr("data-filename");
		var md5 = $(this).attr("data-md5");
		$(this).find(".zip-button").each(function() {
			this.addEventListener('click', function() {
				var modal = $('#zipModal');
				var zipFileInput = $("#zipFileName");
				$("#zipModalForm")[0].reset();
				zipFileInput.attr("value",fileName+".zip");
				modal.openModal();
			}, false);
		});

		$(this).find(".del-button").each(function() {
			this.addEventListener('click', function() {
				var modal = $('#delModal');
				modal.openModal();
			}, false);
		});

		$(this).find(".mov-button").each(function() {
			this.addEventListener('click', function() {
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



