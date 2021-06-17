const Stream = ({ streamName, streamState }) => {
	return (
		<div className={"recording " + (streamState ? "live" : "")}>
			<b>{streamName}</b>
			<p>{streamState ? "Live" : "Not Live"}</p>
		</div>
	);
};

export default Stream;
