export interface ITunesMusic {
  wrapperType: string;
  artworkUrl100: string;
  artistName: string;
  collectionCensoredName: string;
  trackViewUrl: string;
  collectionId: number;
}
export interface ITunesResult {
  resultCount: number;
  results: ITunesMusic[];
}

export const getMusicByName = async (name = ""): Promise<ITunesResult> => {
  return fetch(`https://itunes.apple.com/search?term=${name}`).then(
    (response) => response.json()
  );
};

export const getAlbumById = async (
  id: number | string
): Promise<ITunesResult> => {
  return fetch(`https://itunes.apple.com/lookup?id=${id}`).then(
    (response) => response.json()
  );
};