## osmnotes2csv
### Fetches OpenStreetMap (OSM) notes and saves them as (rfc-4180) csv file.

    Usage:
      ./osmnotes2csv <-bbox=lon,lat,lon,lat> [-limit=n] [-closed] filename
    
    Example:
      ./osmnotes2csv -bbox=7.47,51.84,7.78,52.06 osmnotes.csv
    
    Options:
      -bbox string
        	bounding box (left,bottom,right,top) (required)
      -closed
        	include closed notes
      -limit int
        	maximum number of notes (default 999)
    
    Arguments:
      filename string
            name of csv output file (required)
    
    Remarks:
      A proxy server can be configured via the program environment:
      Temporary: env HTTP_PROXY=http://proxy.server:port ./osmnotes2csv ...
      Permanent: export HTTP_PROXY=http://proxy.server:port
      Permanent: export HTTP_PROXY=http://user:password@proxy.server:port

The prebuild binaries [(Linux, MacOS, Windows, ...)](http://Missing.link) have no dependencies. Just copy and run.

