import React from "react";
import { Card, CardContent, CardMedia, Typography, Paper, Box } from "@mui/material";

interface Artist {
  id: number;
  name: string;
}

interface Format {
  id: number;
  name: string;
}

interface Release {
  id: number;
  title: string;
  year: number;
  artists: Artist[];
  catalogNo: string;
  thumb: string;
  resourceURL: string;
  formats: Format[];
  status: string;
}

interface ReleaseListProps {
  releases: Release[];
}

const ReleaseList: React.FC<ReleaseListProps> = ({ releases }) => {
  return (
    <Box display="flex" flexWrap="wrap" justifyContent="center" gap={2}>
      {releases.map((release, index) => (
        <Paper key={`${release.id}-${index}`} elevation={3} sx={{ width: 300, margin: 2 }}>
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
                Artists: {release.artists?.map((artist) => artist.name).join(", ")}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Year: {release.year}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Formats: {release.formats?.map((format) => format.name).join(", ")}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Catalog No: {release.catalogNo}
              </Typography>

            </CardContent>
          </Card>
        </Paper>
      ))}
    </Box>
  );
};

export default ReleaseList;