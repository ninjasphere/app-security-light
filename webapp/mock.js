function mock() {
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

    var lights = [
      {
        id: "1",
        name: "Front Door Security Light",
        sensors: ["s1"],
        lights: ["d1-1", "d1-2"],
        timeout: 10
      },
      {
        id: "2",
        name: "Bathroom Night Light",
        sensors: ["s2-1", "s2-2"],
        lights: ["d2"],
        timeout: 20,
        timeStart: "20:00",
        timeEnd: "sunrise"
      }
    ];

    var id = lights.length;

    $.mockjax({
      type: "GET",
      url: "/api/security-lights",
      response: function(request) {
        this.responseText = lights;
      }
    });

    $.mockjax({
      type: "POST",
      url: "/api/security-lights",
      response: function(request) {
        console.log("MOCK SAVING", request.data)
        var data = request.data

        if (!data.id) {
          data.id = (id++)+''
        }

        lights = lights.filter(function(l) {
          return l.id != data.id;
        });

        lights.push(data);
        this.responseText = {"error":null};
      }
    });

    $.mockjax({
      type: "DELETE",
      url: /\/api\/security-lights\/(.+)/,
      urlParams: ["id"],
      response: function(request) {
        console.log("deleting", request);

        id = request.urlParams.id;

        lights = lights.filter(function(l) {
          return l.id != id;
        });

        this.responseText = {"error":null};
      }
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
