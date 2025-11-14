import type { FC } from "react";
import { Card, Button } from "react-bootstrap";
import "./MusicCard.css";
import image from "../../assets/DefaultImage.jpg";

interface Props {
  artworkUrl100: string;
  artistName: string;
  collectionCensoredName: string;
  trackViewUrl: string;
  onImageClick?: () => void;
}

export const MusicCard: FC<Props> = ({
  artworkUrl100,
  artistName,
  collectionCensoredName,
  trackViewUrl,
  onImageClick,
}) => (
  <Card className="card">
    <Card.Img
      className="cardImage"
      variant="top"
      src={artworkUrl100 || image}
      height={100}
      width={100}
      onClick={onImageClick}
      style={{ cursor: onImageClick ? "pointer" : "default" }}
    />
    <Card.Body>
      <div className="textStyle">
        <Card.Title>{artistName}</Card.Title>
      </div>
      <div className="textStyle">
        <Card.Text>{collectionCensoredName}</Card.Text>
      </div>
      <Button href={trackViewUrl} target="_blank" variant="primary">
        Открыть в ITunes
      </Button>
    </Card.Body>
  </Card>
);
