import type { FC } from "react";
import type { ITunesMusic } from "../../modules/types";
import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { ALBUMS_MOCK } from "../../modules/mock";
import { Image } from "react-bootstrap";
import defaultImage from "../../assets/DefaultImage.jpg";

export const AlbumPage: FC = () => {
  const [pageData, setPageData] = useState<ITunesMusic | undefined>();
  const { id } = useParams<{ id: string }>();

  useEffect(() => {
    if (!id) return;
    setPageData(
      ALBUMS_MOCK.results.find((a: ITunesMusic) => String(a.collectionId) === id)
    );
  }, [id]);

  if (!pageData) return <div>Album not found</div>;

  return (
    <div>
      <h1>{pageData.collectionCensoredName}</h1>
      <p>Artist: {pageData.artistName}</p>
      <Image
        src={pageData.artworkUrl100 || defaultImage}
        height={200}
        width={200}
      />
    </div>
  );
};
