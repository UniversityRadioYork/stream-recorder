const Recording = ({ recording }) => {
	return (
		<div className="recording">
			<b>{recording.streamName}</b>
			<p>{recording.startTime}</p>
			<a href={"/" + recording.filename}>Download</a>
		</div>
	);
};

export default Recording;
