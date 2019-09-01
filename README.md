# itunes.podcasts

ITunes podcast parsing CLI app. `itupod` provides information about popular podcasts for specific country.
You can get access to the list of genres, show and RSS feed details.

Due to ITunes lookup API [request rate limitation](https://developer.apple.com/library/archive/documentation/AudioVideo/Conceptual/iTuneSearchAPI/Searching.html#//apple_ref/doc/uid/TP40017632-CH5-SW1) (approximately 20 calls per minute) loading of the show details can take some time.
App provides a way to load details by chunks with proper delay between API request (see details below).

## How to use it

Here is the list of possible commands for retrieve data from ITunes:

- `itupod [-g | -genre]` - this will load list of genres and save in current folder
- `itupod [-s | -show] PATH_TO_GENRES` - this will load list of shows and save in current folder. You must specify a path to `genres.json` file in arguments
- `itupod [-d | -details] [-chunk] PATH_TO_SHOWS` - this will load chunk sized list of show details and save in current folder. You must specify a path to `shows.json` file in arguments
- `itupod [-f | -feed] PATH_TO_DETAILS` - this will load feed and save in current folder. You must specify a path to `shows.details.json` file in arguments

By default files will be stored into the `/tmp` folder, you can change it be providing `-out` flag with path for desired folder
