import type { FC } from "react";
import type { ITunesMusic } from "../../modules/types";
import { useState } from "react";
import { Col, Row, Spinner } from "react-bootstrap";
import { MusicCard } from "../../components/MusicCard/MusicCard";
import { ALBUMS_MOCK } from "../../modules/mock";

const ITunesPage: FC = () => {
  const [searchValue, setSearchValue] = useState<string>("");
  const [loading, setLoading] = useState<boolean>(false);
  const [music, setMusic] = useState<ITunesMusic[]>([]);

  const handleSearch = () => {
    setLoading(true);
    // Здесь обычно API вызов, но используем mock
    setMusic(
      ALBUMS_MOCK.results.filter((item: ITunesMusic) =>
        item.collectionCensoredName
          .toLocaleLowerCase()
          .startsWith(searchValue.toLocaleLowerCase())
      )
    );
    setLoading(false);
  };

  return (
    <div>
      <input
        value={searchValue}
        onChange={(e) => setSearchValue(e.target.value)}
        placeholder="Search"
      />
      <button onClick={handleSearch}>Search</button>
      {loading && <Spinner animation="border" />}
      <Row>
        {music.map((item: ITunesMusic) => (
          <Col key={item.collectionId}>
            <MusicCard
              artworkUrl100={item.artworkUrl100}
              artistName={item.artistName}
              collectionCensoredName={item.collectionCensoredName}
              trackViewUrl={item.trackViewUrl}
            />
          </Col>
        ))}
      </Row>
    </div>
  );
};

export default ITunesPage;
