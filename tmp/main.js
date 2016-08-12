// Log date = X axis
// Log tag = series name
// Log message = datapoint

var chart;
var lastx = 0;
var lasty = 0;

var Chart = function(containerId, title) {
    this.cjsChart = new CanvasJS.Chart(containerId, {
        title: {
            text: title
        },
        legend: {
			horizontalAlign: "left",
			verticalAlign: "center"
		},
		axisY:{
			includeZero: false,
		},
		data: []
    });
    this.seriesIds = {}
}

Chart.prototype.addSeries = function(seriesname, seriesOptions) {
    var sid = this.cjsChart.options.data.length;
    this.seriesId[seriesname] = sid;
    this.cjsChart.options.data.push(seriesOptions);
    this.cjsChart.options.data[sid].dataPoints |= []
}

Chart.prototype.addPoint = function(seriesname, data) {
    if !seriesname || !datapoint {
        return
    }

    if this.seriesIds[seriesname] == undefined {
        this.addSeries(seriesname, data.seriesOptions);
    }

    var sid = this.seriesIds[seriesname];
    this.cjsChart.options.data[sid].dataPoints.push(data.pointData)
}

Chart.prototype.render = function() {
    this.cjsChart.render();
}

Chart.gc = function() {
    series = this.cjsChart.options.data;
    for (var i in series) {
        var s = series[i]
        var l = series[i].dataPoints.length
        var chartMin = this.cjsChart.options.axisX.minimum;
        if s.dataPoints[l/2].x < chartMin {
            // TODO slice that bad boy in half
        }
    }
}

Chart.clearData = function(seriesname) {
    this.seriesId = {};
    this.cjsChart.options.data = []
}

window.onload = function() {
    chart = new Chart("chartContainer", "Chart title");
    chart.addSeries("Series 1", {})
    setInterval(chart.gc, 1000)
}

function doUpdate() {
    chart.addPoint("Series 1", {x: lastx, y: lasty})
    chart.render()
    lasty += (Math.random() * 10 - 5);
    lastx += 1
    console.log("updating!")
}
//
//function doUpdate0(chart) {
//    lasty += (Math.random() * 10 - 5);
////    if lastx % 10 == 0 {
////        chart.options.data[1].dataPoints.push({x: lastx, y: lasty+ (Math.random() * 30 - 15), z: (Math.random()*50)})
////        chart.options.data[1].dataPoints.shift()
////    } else {
//        chart.options.data[0].dataPoints.push({x: lastx, y: lasty})
//        chart.options.data[0].dataPoints.shift()
////    }
//    chart.render()
//    lastx += 1
//    console.log("updating!")
//}
//
//var limit = 100;    //increase number of dataPoints by increasing this
//var data = [];
//var dataSeries = { type: "line" };
//var dataPoints = [];
//for (var i = 0; i < limit; i += 1) {
//    lasty += (Math.random() * 10 - 5);
//    dataPoints.push({
//        x: i - limit / 2,
//        y: lasty
//    });
//}
//lastx = limit/2
//dataSeries.dataPoints = dataPoints;
//data.push(dataSeries);
//
//data.push({
//		type: "bubble",
//		//markerType: "square",
//		showInLegend: true,
//		dataPoints: [
//			{ x: 10, y: 71, z:12 },
//			{ x: 20, y: 55, z:23 },
//			{ x: 30, y: 50, z:6  },
//			{ x: 40, y: 65, z:2  },
//			{ x: 50, y: 95, z:50 },
//			{ x: 60, y: 68, z:13 },
//			{ x: 70, y: 28, z:21 },
//			{ x: 80, y: 34, z:5  },
//			{ x: 100, y: 14, z:7  }
//		]
//		});
//
//
//setInterval(doUpdate, 100)

