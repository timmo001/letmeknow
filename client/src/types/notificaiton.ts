export type Image = {
  height?: number | null;
  width?: number | null;
  url: string;
};

export type Notification = {
  title?: string | null;
  subtitle?: string | null;
  content?: string | null;
  image?: Image | null;
};
