import React, { useEffect, useState } from 'react';
import { fetchRelations } from '../services/api';
import { Typography, CircularProgress, Box, List, ListItem, ListItemText } from '@mui/material';

interface Relation {
  tableName: string;
  relations: string[];
}

const TablesRelations: React.FC = () => {
  const [relations, setRelations] = useState<Relation[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    const loadRelations = async () => {
      try {
        const data = await fetchRelations();
        const formattedRelations = Object.keys(data).map(tableName => ({
          tableName,
          relations: data[tableName],
        }));
        setRelations(formattedRelations);
      } catch (err: any) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    loadRelations();
  }, []);

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" height="100vh">
        <CircularProgress />
      </Box>
    );
  }

  if (error) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" height="100vh">
        <Typography color="error">{error}</Typography>
      </Box>
    );
  }

  return (
    <Box my={4}>
      <Typography variant="h4" component="h1" gutterBottom>
        Tables and Relations
      </Typography>
      {relations.map(({ tableName, relations }) => (
        <Box key={tableName} mb={3}>
          <Typography variant="h5" component="h2">
            {tableName}
          </Typography>
          <List>
            {relations.map((relation, index) => (
              <ListItem key={index}>
                <ListItemText primary={relation} />
              </ListItem>
            ))}
          </List>
        </Box>
      ))}
    </Box>
  );
};

export default TablesRelations;