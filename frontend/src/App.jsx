import "./App.css";
import { useState, useEffect } from "react";
import Recording from "./Recording";
import Stream from "./Stream";

function App() {
	const [recordings, setRecordings] = useState([]);
	const [streams, setStreams] = useState({});

	useEffect(() => {
		let endpoint =
			process.env.NODE_ENV === "production"
				? "/recordings-json"
				: "http://localhost:3001/recordings-json";
		fetch(endpoint)
			.then((res) => res.json())
			.then((data) => setRecordings(data));

		let ws = new WebSocket("ws://localhost:3001/ws");
		ws.onmessage = (ev) => {
			let message = ev.data.split(":");
			setStreams((s) => ({
				...s,
				[message[0]]: message[1] === "true",
			}));
		};
		ws.onopen = () => ws.send("QUERY");
		ws.onclose = () => window.location.reload;
	}, []);

	return (
		<div className="App">
			<h1>((URY)) Stream Recorder</h1>
			<h2> Live </h2>
			<div className="flex-container">
				{Object.keys(streams).map((stream, idx) => (
					<Stream
						key={"stream" + idx}
						streamName={stream}
						streamState={streams[stream]}
					/>
				))}
			</div>
			<h2> Recordings </h2>
			<div className="flex-container">
				{recordings.map((recording, idx) => (
					<Recording key={"recording" + idx} recording={recording} />
				))}
			</div>
		</div>
	);
}

export default App;
