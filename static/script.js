// I hate JavaScript
//               -- Howl

if (!String.prototype.format) {
  String.prototype.format = function () {
    var args = arguments;
    return this.replace(/{(\d+)}/g, function (match, number) {
      return typeof args[number] != 'undefined'
        ? args[number]
        : match
				;
    });
  };
}

var printResult = function (jqObj, scoreID) {
	if (typeof scoreID === "undefined") {
		return;
	}
	$.getJSON("/api/v1/score?id=" + scoreID, function (resp) {
		if (!resp.ok) {
			$(jqObj).html('<span class="mui--text-accent">Something went wrong with this file! ({0})</span>'.format(resp.message));
			return;
		}
		if (resp.calculated > 1) {
			$(jqObj).html('<span class="mui--text-accent">Something went wrong with this file! Looks like a server error. You may want to tell the website owner about this.</span>');
			return;
		}
		if (resp.calculated == 0) {
			resp.score.pp = "Hold on, still calculating...";
		} else {
			resp.score.pp += "pp";
		}
		var modsStr = resp.score.mods_str != "" ? " +" + resp.score.mods_str : "";
		$(jqObj).html("<h1><b>{0}</b></h1><h4 class='mui--text-dark-secondary'>{1} on {2} - {3} [{4}] ({5}){6}</h4>".format(resp.score.pp, resp.score.player, resp.beatmap.author, resp.beatmap.title, resp.beatmap.diff_name, resp.beatmap.creator, modsStr));
		if (resp.calculated == 0) {
			window.setTimeout(function () {
				printResult(jqObj, scoreID);
			}, 5000);
		}
	});
};

$(function () {
	'use strict';
	var url = '/api/v1/score';
	$('#file').fileupload({
		url: url,
		dataType: 'json',
		autoUpload: true,
		acceptFileTypes: /(\.)(osr|osu)$/i,
		maxFileSize: 1024 * 1024,
	}).on('fileuploadadd', function (e, data) {
		data.context = $('<div/>').addClass("mui-panel").appendTo('#cont');
		var that = $(this).data('blueimpFileupload'); 
		$.each(data.files, function (index, file) {
			var ext = file.name.split(".").pop();
			if (ext == "osu") {
				data.paramName = "beatmap";
				that.options.url = "/api/v1/beatmap";
			} else if (ext != "osr") {
				// I haven't got the slightest idea of how this works. But it does. So yeah.
				$(data.context).html($('<span class="mui--text-accent"/>').text("Please upload either an .osu file or .osr file."));
				e.preventDefault();
				data.files = [];
				return false;
			} else {
				that.options.url = "/api/v1/score";				
			}
			var node = $('<p/>')
				.append($('<span/>').text("Uploading " + file.name + "..."));
			node.appendTo(data.context);
		});
	}).on('fileuploaddone', function (e, data) {
		if (data.url == "/api/v1/beatmap" && data.result.ok) {
			$(data.context).append(" done!");
		}
		if (data.result.ok && data.result.score_id != 0) {
			printResult(data.context, data.result.score_id);
		} else if (file.error) {
		    var error = $('<span class="mui--text-accent"/>').text(file.message);
		    $(data.context.children()[index])
				.append(error);
		}
	}).on('fileuploadfail', function (e, data) {
		// hackish but who cares
		var errorMessage = data.response().jqXHR.responseJSON.message;
		var error = $('<span class="mui--text-accent"/>').text(errorMessage);
		if (errorMessage) {
			$(data.context)
				.append(error);
		}
	}).prop('disabled', !$.support.fileInput)
		.parent().addClass($.support.fileInput ? undefined : 'disabled');
});
$("#label-click").click(function () {
	$("#file").click();
});
