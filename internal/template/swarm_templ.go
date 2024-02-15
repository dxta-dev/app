// Code generated by templ - DO NOT EDIT.

// templ: version: 0.2.476
package template

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import (
	"github.com/dxta-dev/app/internal/data"

	"time"
)

type SwarmProps struct {
	Series               SwarmSeries
	StartOfTheWeek       time.Time
	EventIds             []int64
	EventMergeRequestIds []int64
}

var eventTypeNames = map[data.EventType]string{
	data.UNKNOWN:                "UNKNOWN",
	data.OPENED:                 "OPENED",
	data.STARTED_CODING:         "STARTED_CODING",
	data.STARTED_PICKUP:         "STARTED_PICKUP",
	data.STARTED_REVIEW:         "STARTED_REVIEW",
	data.NOTED:                  "NOTED",
	data.ASSIGNED:               "ASSIGNED",
	data.CLOSED:                 "CLOSED",
	data.COMMENTED:              "COMMENTED",
	data.COMMITTED:              "COMMITTED",
	data.CONVERT_TO_DRAFT:       "CONVERT_TO_DRAFT",
	data.MERGED:                 "MERGED",
	data.READY_FOR_REVIEW:       "READY_FOR_REVIEW",
	data.REVIEW_REQUEST_REMOVED: "REVIEW_REQUEST_REMOVED",
	data.REVIEW_REQUESTED:       "REVIEW_REQUESTED",
	data.REVIEWED:               "REVIEWED",
	data.UNASSIGNED:             "UNASSIGNED",
}

func getEventName(eventType data.EventType) string {
	if name, exists := eventTypeNames[eventType]; exists {
		return name
	}
	return "Invalid EventType"
}

func getChart(chartId string, endpoint string, circleIds []int64, circleMergeRequestIds []int64) templ.ComponentScript {
	return templ.ComponentScript{
		Name: `__templ_getChart_cb35`,
		Function: `function __templ_getChart_cb35(chartId, endpoint, circleIds, circleMergeRequestIds){if (circleIds === null) {
		return;
	}

	const svg = document.querySelector(` + "`" + `${chartId} > svg` + "`" + `);
    const circles = document.querySelectorAll(` + "`" + `${chartId} > svg > circle` + "`" + `);

    function moveToTop(element) {
        element.parentNode.appendChild(element);
    }

	function getCircleInfo(mrid) {
		const searchParams = new URLSearchParams(document.location.search);
		searchParams.set('mr', mrid);
		const endpointWithMrid = ` + "`" + `${endpoint}${mrid}` + "`" + `;
		const mrUrl = new URL(endpointWithMrid, document.location.origin);
		htmx.ajax('GET', mrUrl.toString(), '#slide-over');
	}

	function createLineBetweenCircles(x1, y1, x2, y2) {
		const line = document.createElementNS("http://www.w3.org/2000/svg", "line");

		line.setAttribute("x1", x1);
		line.setAttribute("y1", y1);
		line.setAttribute("x2", x2);
		line.setAttribute("y2", y2);
		line.setAttribute("stroke", "blue");
		line.setAttribute("stroke-width", "2");

		return line;
	}

	function calculateDistance(x1, y1, x2, y2) {
		return Math.pow(x2 - x1, 2) + Math.pow(y2 - y1, 2);
	}

	function orderCircles(circles) {
		if (circles.length === 0) return [];
		if (circles.length === 1) return circles;

		let unorderedCircles = [...circles];

		let currentCircle = unorderedCircles.shift();

		const orderedCircles = [currentCircle];

		while (unorderedCircles.length > 0) {
			let closestCircleIndex = 0;
			let minDistance = Infinity;

			for (let i = 0; i < unorderedCircles.length; i++) {
				const distance = calculateDistance(currentCircle.getAttribute("cx"), currentCircle.getAttribute("cy"), unorderedCircles[i].getAttribute("cx"), unorderedCircles[i].getAttribute("cy"));
				if (distance < minDistance) {
					minDistance = distance;
					closestCircleIndex = i;
				}
			}

			let closestCircle = unorderedCircles.splice(closestCircleIndex, 1)[0];

			currentCircle = closestCircle;

			orderedCircles.push(currentCircle);
		}

		return orderedCircles;
	}


	const searchParams = new URLSearchParams(document.location.search);

	for (let i = 0; i < circles.length; i++) {
		circles[i].setAttribute("data-id", circleIds[i]);
		circles[i].setAttribute("data-merge-request-id", circleMergeRequestIds[i]);
	}

	let size = +circles[0].getAttribute("r");

	circles.forEach((circle, i) => {
		let strokes = [];
		let lines = [];
		circle.addEventListener("mouseover", (e) => {
			const mergeRequestId = circle.getAttribute("data-merge-request-id");
			let circles = document.querySelectorAll(` + "`" + `${chartId} > svg > circle[data-merge-request-id="${mergeRequestId}"]` + "`" + `);
			circles = orderCircles(circles);
			prevCircle = circles[0];
			for (let i = 0; i < circles.length; i++) {
				circles[i].setAttribute("r", size + 2);
				strokes.push(circles[i].style.stroke);
				circles[i].style.stroke = "blue";
				circles[i].style.strokeWidth = "2";
				if (i > 0 && 300 < calculateDistance(prevCircle.getAttribute("cx"), prevCircle.getAttribute("cy"), circles[i].getAttribute("cx"), circles[i].getAttribute("cy"))){
					const line = createLineBetweenCircles(prevCircle.getAttribute("cx"), prevCircle.getAttribute("cy"), circles[i].getAttribute("cx"), circles[i].getAttribute("cy"));
					lines.push(line);
					svg.appendChild(line);
					prevCircle = circles[i];
				}
			}
			for (let i = 0; i < circles.length; i++) {
				moveToTop(circles[i]);
			}
		});

		circle.addEventListener("mouseout", (e) => {
			const mergeRequestId = circle.getAttribute("data-merge-request-id");
			const circles = document.querySelectorAll(` + "`" + `${chartId} > svg > circle[data-merge-request-id="${mergeRequestId}"]` + "`" + `);
			for (let i = 0; i < circles.length; i++) {
				circles[i].setAttribute("r", size);
				circles[i].style.stroke = strokes[i];
				circles[i].style.strokeWidth = "1";
			}
			for (let i = 0; i < lines.length; i++) {
				lines[i].remove();
			}
			strokes = [];
			lines = [];
		});

		circle.addEventListener("click", (e) => {
			getCircleInfo(Number(circle.getAttribute("data-merge-request-id")));

		});
	});}`,
		Call:       templ.SafeScript(`__templ_getChart_cb35`, chartId, endpoint, circleIds, circleMergeRequestIds),
		CallInline: templ.SafeScriptInline(`__templ_getChart_cb35`, chartId, endpoint, circleIds, circleMergeRequestIds),
	}
}

func swarm(props SwarmProps) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"flex items-center justify-center\" id=\"swarm-chart\"><style text=\"text/css\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Var2 := `
			svg > circle:hover {
				cursor: pointer;
			}
		`
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ_7745c5c3_Var2)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</style>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = swarmChartComponent(props.Series, props.StartOfTheWeek).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = getChart("#swarm-chart", "merge-request/", props.EventIds, props.EventMergeRequestIds).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}