import "./App.css";
import { useState, useEffect } from "react";
import Recording from "./Recording";

function App() {
	const [recordings, setRecordings] = useState([]);

	useEffect(() => {
		fetch("/recordings-json")
			.then((res) => res.json())
			.then((data) => setRecordings(data));
	}, []);

	return (
		<div className="App">
			<h1> Live </h1>
			<div className="flex-container"></div>
			<br />
			<h1> Recordings </h1>
			<div className="flex-container">
				{recordings.map((recording, idx) => (
					<Recording key={idx} recording={recording} />
				))}
			</div>
		</div>
	);
}

export default App;
