import { h } from 'preact';
import { useEffect, useRef, useState } from 'preact/hooks';
import wasmLib from '../../lib/wasm';
import style from './style.css';
import * as d3 from "d3";
import getGraphData from '../../lib/graph';
/**
 * Ulitlity function that takes a  random string and returns a HEX color
 * Cfr. https://stackoverflow.com/questions/3426404/create-a-hexadecimal-colour-based-on-a-string-with-javascript
 * 
 * @param {*} str 
 * @returns 
 */
const stringToColour = (str) => {
    var hash = 0;
    for (var i = 0; i < str.length; i++) {
        hash = str.charCodeAt(i) + ((hash << 5) - hash);
    }
    var colour = "#";
    for (var i = 0; i < 3; i++) {
        var value = (hash >> (i * 8)) & 0xff;
        colour += ("00" + value.toString(16)).substr(-2);
    }
    return colour;
};


const Graph = ({ notionKey }) => {
    const divRef = useRef(null);

    const [data, setData] = useState();

    if (!data) {
        // Getdata is a wasm function
        getGraphData(notionKey).then(res => {
            setData(res);
        });
    }

    const margin = { top: 10, right: 30, bottom: 30, left: 40 };

    useEffect(() => {

        if (!data) {
            return;
        }

        function handleZoom(e) {
            // apply transform to the chart
            d3.select(divRef.current).select("svg g").attr("transform", e.transform);
        }

        let zoom = d3.zoom().on("zoom", handleZoom);

        // append the svg object to the body of the page
        var svg = d3
            .select(divRef.current)
            .append("svg")
            .attr("width", "100%")
            .attr("height", "100%")
            .append("g")
            .attr("transform", "translate(" + margin.left + "," + margin.top + ")");


        const dataNodes = data.Pages.map((page) => {
            return {
                id: page.Id,
                name: page.Id,
                color: stringToColour(page.ParentDatabaseId),
            };
        });

        const dataLinks = data.Links.map((link) => {
            return {
                source: link.FromPage,
                target: link.ToPage,
            };
        });
        // Initialize the links

        var link = svg
            .selectAll("line")
            .data(dataLinks)
            .enter()
            .append("line")
            .style("stroke", "#aaa");

        // Initialize the nodes
        var node = svg
            .selectAll("circle")
            .data(dataNodes)
            .enter()
            .append("circle")
            .attr("r", 10)
            .style("fill", "#69b3a2");

        // Let's list the force we wanna apply on the network
        var simulation = d3
            .forceSimulation(dataNodes) // Force algorithm is applied to data.nodes
            .force(
                "link",
                d3
                    .forceLink() // This force provides links between nodes
                    .id(function (d) {
                        return d.id;
                    }) // This provide  the id of a node
                    .links(dataLinks) // and this the list of links
            )
            .force("charge", d3.forceManyBody().strength(-300)) // This adds repulsion between nodes. Play with the -400 for the repulsion strength
            .force("center", d3.forceCenter()) // This force attracts nodes to the center of the svg area
            .on("end", ticked);

        // This function is run at each iteration of the force algorithm, updating the nodes position.
        function ticked() {
            link
                .attr("x1", (d) => d.source.x)
                .attr("y1", (d) => d.source.y)
                .attr("x2", (d) => d.target.x)
                .attr("y2", (d) => d.target.y);

            node.attr("cx", (d) => d.x).attr("cy", (d) => d.y);

            node.style("fill", (d) => d.color);
        }

        d3.select(divRef.current).select("svg").call(zoom);


    })

    return (<div ref={divRef} class={style.graph}></div>);

}





export default Graph;



// const go = new Go();
// WebAssembly.instantiateStreaming(
//     fetch("main.wasm"),
//     go.importObject
// ).then((result) => {
//     go.run(result.instance);
// });



// set the dimensions and margins of the grap