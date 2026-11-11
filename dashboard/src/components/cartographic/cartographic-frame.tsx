export function CartographicFrame() {
  return (
    <div className="cartographic-frame" aria-hidden="true">
      <div className="longitude longitude-one">67°30′ O</div>
      <div className="longitude longitude-two">67°00′ O</div>
      <div className="longitude longitude-three">66°30′ O</div>
      <div className="latitude latitude-one">16°00′ S</div>
      <div className="latitude latitude-two">16°30′ S</div>
      <div className="latitude latitude-three">17°00′ S</div>
      <span className="survey-cross cross-one" />
      <span className="survey-cross cross-two" />
      <span className="survey-cross cross-three" />
      <div className="map-scale"><span>0</span><span>5</span><span>10</span><span>20 km</span></div>
    </div>
  );
}
