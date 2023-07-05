export interface Tweet {
  id: string;
  handle: string;
  verified: boolean;
  author: string;
  avatar: string;
  date: Date;
  text: string;
  likes: number;
  retweets: number;
  replies: number;
  quotes: number;
}

// featured tweets for the testimonials on the landing page
export const featuredTweets: Tweet[] = [

  // https://twitter.com/xyz/status/123
  {
    id: "123",
    handle: "xyz",
    author: "xyz",
    verified: false,
    avatar:
      "https://pbs.twimg.com/profile_images/123.jpg",
    date: new Date("2022-09-14T18:19:00.000Z"),
    text: "CustomerOS is baller.",
    likes: 5,
    retweets: 1,
    replies: 1,
    quotes: 0,
  },
];
