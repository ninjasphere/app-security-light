$(function() {
  $('.add').click(function() {
    edit({
      name: "My Security Light",
      sensors: [],
      lights: [],
      timeout: 5,
      time: null
    })
  });

  $('.save').click(function() {
    // TODO: Save the thing

    list();
  });

  $(document).on('click', 'button.edit', function(e){
    edit(securityLights[$(e.target).data('id')]);
  })
})

var securityLights = {};

function list() {

  $('#edit').hide()
  $('#list').show()

  $('#securityLights').empty();
  $.get("/api/security-lights", function(l) {
    l.forEach(function(light) {
      securityLights[light.id] = light;
      $('#securityLights').append('<li>' + light.name + ' <button class="edit" data-id="' + light.id + '">Edit</button> <button class="delete" data="' + light.id + '">Delete</button></li>');
    })
  });

}

function edit(securityLight) {

  function fill(el, url, selected) {
    el.empty();
    $.get(url, function(items) {
      items.forEach(function(item) {
        el.append('<option value="' + item.id + '"' + (selected.indexOf(item.id) > -1?' selected':'') + '>' + item.name + item.id + '</option>');
      })
    })
  }

  fill($('[name=sensors]'), "/api/sensors", securityLight.sensors)
  fill($('[name=lights]'), "/api/lights", securityLight.lights)

  $('[name=id]').val(securityLight.id || "")
  $('[name=name]').val(securityLight.name || "")
  $('[name=timeout]').val(securityLight.timeout || "5")
  if (securityLight.time) {
      $('[name=active][value=between]').click()
      $('[name=timeStart]').val(securityLight.time.start)
      $('[name=timeEnd]').val(securityLight.time.end)
  } else {
      $('[name=active][value=always]').click()
      $('[name=timeStart]').val('sunrise')
      $('[name=timeEnd]').val('sunset')
  }

  $('#list').hide()
  $('#edit').show()
}

// Times dropdowns
function pad(a,b){return(1e15+a+"").slice(-b)}

var times = ['midnight', 'sunrise', 'dawn', 'midday', 'sunset', 'dusk']

for (var h = 0; h < 24; h++) {
  for (var m = 0; m < 60; m += 15) {
    times.push(pad(h, 2) + ':' + pad(m, 2));
  }
}
$(function(){
  var els = $('.times');
  times.forEach(function(t) {
    els.append('<option>'+t+'</option>');
  });

  /* Don't work right with 'sunrise' etc.
  var endOptions = $('[name=timeEnd] option');

  $('[name=timeStart]').change(function(e) {
    console.log("Change", e);
    var selected = $(e.srcElement).find('option').not(function(){ return !this.selected })[0].index;

    endOptions.forEach(function(el, index) {
      if (index > selected) {
        $(el).removeAttr('disabled');
      } else {
        $(el).attr('disabled', 'disabled');
      }
    });

    if (endOptions.not(function(){ return !this.selected })[0].disabled) {
      endOptions.not(function(){ return this.disabled})[0].selected = true
    }
  });
  */
})






// Http mocking
if (location.protocol != "file:") {

  list();

} else {

  ;(function ($) {
      $.getScript = function(src, func) {
          var script = document.createElement('script');
          script.async = "async";
          script.src = src;
          if (func) {
             script.onload = func;
          }
          document.getElementsByTagName("head")[0].appendChild( script );
      }
  })($)

  $.getScript("zepto-mockjax.js", function() {

    $.mockjax({
      type: "GET",
      url: "/api/security-lights",
      responseText: [
        {
          id: "1",
          name: "Front Door Security Light",
          sensors: ["s1"],
          lights: ["d1-1", "d1-2"],
          timeout: 10,
          time: null
        },
        {
          id: "2",
          name: "Bathroom Night Light",
          sensors: ["s2-1", "s2-2"],
          lights: ["d2"],
          timeout: 20,
          time: {
            start: "20:00",
            end: "sunrise"
          }
        }
      ]
    });

    $.mockjax({
      type: "GET",
      url: "/api/lights",
      responseText: [
        {
          id: "d1-1",
          name: "Front Door Overhead"
        },
        {
          id: "d1-2",
          name: "Front Door Spotlight"
        },
        {
          id: "d2",
          name: "Bathroom Lamp"
        }
      ]
    });

    $.mockjax({
      type: "GET",
      url: "/api/sensors",
      responseText: [
        {
          id: "s1",
          name: "Front Door Motion"
        },
        {
          id: "s2-1",
          name: "Bathroom Motion"
        },
        {
          id: "s2-2",
          name: "Bathroom Door Motion"
        },
      ]
    });

    list();

  });
}
