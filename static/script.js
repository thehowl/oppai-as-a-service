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
var entityMap = {
	"&": "&amp;",
	"<": "&lt;",
	">": "&gt;",
	'"': '&quot;',
	"'": '&#39;',
	"/": '&#x2F;'
};

function escapeHtml(string) {
	return String(string).replace(/[&<>"'\/]/g, function (s) {
		return entityMap[s];
	});
}

function genPanel() {
	return $('<div/>').addClass("mui-panel").appendTo('#cont');
}

var curPath = "/" + (window.location.search == "" ? "?" : window.location.search);

var printResult = function (jqObj, scoreID, alreadyDone) {
	if (typeof scoreID === "undefined") {
		return;
	}
	if (!alreadyDone) {
		curPath += "&p=" + scoreID;
		window.history.pushState("", "", curPath);
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
		$(jqObj).html(
			"<h1><b>{0}</b></h1><h4 class='mui--text-dark-secondary'>{1} on {2} - {3} [{4}] ({5}){6}</h4>"
				.format(
					resp.score.pp, 
					escapeHtml(resp.score.player),
					escapeHtml(resp.beatmap.author),
					escapeHtml(resp.beatmap.title),
					escapeHtml(resp.beatmap.diff_name),
					escapeHtml(resp.beatmap.creator),
					modsStr
				)
		);
		if (resp.calculated == 0) {
			window.setTimeout(function () {
				printResult(jqObj, scoreID, true);
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
		data.context = genPanel();
		var that = $(this).data('blueimpFileupload');
		$.each(data.files, function (index, file) {
			var ext = file.name.split(".").pop();
			if (ext == "osu") {
				data.paramName = "beatmap";
				that.options.url = "/api/v1/beatmap";
			} else if (ext != "osr") {
				// I haven't got the slightest idea of how this works. But it does. So yeah.
				$(data.context).html($('<span class="mui--text-accent"/>').text("Please upload either an .osu file or .osr file."));
				data.context = $("<div hidden />");
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

function donationModal() {
	// initialize modal element
	var modalEl = document.createElement('div');
	modalEl.style.width = '400px';
	modalEl.style.height = '400px';
	modalEl.style.margin = '100px auto';
	modalEl.style.backgroundColor = '#fff';
	modalEl.className = "mui-panel";
	modalEl.innerHTML = "<p>So, someone asked for it, and there we go I guess.</p><p>When you want to make a donation to OaaS, you must think on what you believe of this service to be great.</p><ol><li>The PP formula. In that case, get osu! supporter (Tom94 works for osu!).</li><li>Oppai, the magic program that calculates the PP for any beatmap. In that case, <a href='https://www.paypal.com/cgi-bin/webscr?cmd=_s-xclick&hosted_button_id=5E289LJ5UUG3Q'>click here to donate to its developer</a>.</li><li>Finally, me. The dude who made the web application. If that's the case, just send some money via PayPal to dahhowl@gmail.com. Nope, no fancy link. Can't be bothered.</li></ol>"; 

	// show modal
	mui.overlay('on', modalEl);
}
