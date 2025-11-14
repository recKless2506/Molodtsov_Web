// mock.ts
export interface ITunesMusic {
  wrapperType: string;
  artistName: string;
  collectionCensoredName: string;
  trackViewUrl: string;
  artworkUrl100: string;
  collectionId: number;
}

export const ALBUMS_MOCK = {
  resultCount: 3,
  results: [
    {
      wrapperType: "track",
      artistName: "Pink Floyd",
      collectionCensoredName: "The Wall",
      trackViewUrl: "https://itunes.apple.com/album/the-wall",
      artworkUrl100: "",
      collectionId: 1,
    },
    {
      wrapperType: "track",
      artistName: "Queen",
      collectionCensoredName: "A Night At The Opera",
      trackViewUrl: "https://itunes.apple.com/album/a-night-at-the-opera",
      artworkUrl100: "",
      collectionId: 2,
    },
    {
      wrapperType: "track",
      artistName: "AC/DC",
      collectionCensoredName: "Made in Heaven",
      trackViewUrl: "https://itunes.apple.com/album/made-in-heaven",
      artworkUrl100: "",
      collectionId: 3,
    },
  ],
};
