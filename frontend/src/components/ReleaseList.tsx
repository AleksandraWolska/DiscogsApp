import React from "react";
import { Card, CardContent, CardMedia, Typography, Link, Grid } from "@mui/material";

interface Release {
  id: number;
  title: string;
  year: number;
  artist: string;
  format: string;
  catalogNo: string;
  thumb: string;
  resourceURL: string;
}

interface ReleaseListProps {
  releases: Release[];
}

const ReleaseList: React.FC<ReleaseListProps> = ({ releases }) => {
  return (
    <Grid container spacing={2}>
      {releases.map((release) => (
        <Grid item xs={12} sm={6} md={4} key={release.id}>
          <Card>
            <CardMedia
              component="img"
              height="140"
              image={release.thumb || `${process.env.PUBLIC_URL}/cover.png`}
              alt={release.title}
            />
            <CardContent>
              <Typography gutterBottom variant="h5" component="div">
                {release.title}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Artist: {release.artist}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Year: {release.year}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Format: {release.format}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Catalog No: {release.catalogNo}
              </Typography>
              <Link href={release.resourceURL} target="_blank" rel="noopener">
                More Info
              </Link>
            </CardContent>
          </Card>
        </Grid>
      ))}
    </Grid>
  );
};

export default ReleaseList;