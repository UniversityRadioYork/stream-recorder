import "./App.css";
import { useState, useEffect } from "react";
import { Recording } from "./Recording";

function App() {
	const [recordings, setRecordings] = useState([]);

	useEffect(() => {
		fetch("/recordings")
			.then((res) => res.json())
			.then((data) => setRecordings(data));
	}, []);

	return (
		<div className="App">
			<h1> Live </h1> <h1> Recordings </h1>{" "}
			{recordings.map((recording, idx) => (
				<Recording key={idx} recording={recording} />
			))}{" "}
		</div>
	);
}

export default App;
