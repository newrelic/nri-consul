// $( document ).ready(function() {
//     // $('.resource').hide();
//     // $('.metric').hide();
//     $('.resource ul').hide();
// 	$('.metric-header').hide();
// 	$('#left-col').show();
//     $('.resource_type').find('h3').click(function() {
//     	$(this).parent().find('.resource').slideToggle();
//     	$(this).parent().toggleClass('active');
//     });
//     $('.resource').find('h4').click(function() {
//     	$('.resource').find('h4').css('color', 'inherit');
//     	$('#right-col h4').css('color', '#inherit');
// 		$('#right-col').html($(this).parent().parent().html()).show().find('.metric').show();
// 		$('#right-col').children().addClass('list-group-item');
// 		$('#right-col').find('.event').show();
//     	$('#right-col').find('ul').slideDown();
//     	$(this).css('color', '#5bc0de');
//     	$('#right-col h4').css('color', '#5bc0de');
// 		$('#right-col h4').show().click(function() {
// 		    $(this).parent().find('.resource-content-group').slideToggle();
// 		});
// 		$('.info .resource-content-group').hide();
//     });
// });

$(document).ready(function() {
    $('.entity-type').click(function() {
        $(this).parent().find('.entity-list').slideToggle();
    })
    $('.entity-list').find('div').click(function() {
        $('#entityDetails').find('.entity-details').slideUp()
        var id = "#" + $(this).attr('id') + "Details"
        $(id).slideToggle()
    })
});