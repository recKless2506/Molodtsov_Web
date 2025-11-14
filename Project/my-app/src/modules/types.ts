export interface ITunesMusic {
  wrapperType: string;
  artistName: string;
  collectionCensoredName: string;
  trackViewUrl: string;
  artworkUrl100: string;
  collectionId: number;
}

export interface ITunesResult {
  resultCount: number;
  results: ITunesMusic[];
}
