const Recording = ({ recording }) => {
	const d = new Date(recording.startTime);

	return (
		<div className="recording">
			<b>{recording.streamName}</b>
			<p>{d.toLocaleString()}</p>
			<a href={"/" + recording.filename}>Download</a>
		</div>
	);
};

export default Recording;
