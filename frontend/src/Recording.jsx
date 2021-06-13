const Recording = ({ recording }) => {
	return (
		<div>
			<b>{recording.streamName}</b>
			<p>{recording.startTime}</p>
			<a href={"/recordings/" + recording.filename}>Download</a>
		</div>
	);
};

export default Recording;
