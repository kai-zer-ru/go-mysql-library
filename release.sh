#!/usr/bin/env bash
releaseMsg=$1
version=$2
curl 'https://gitlab.sarafann.com/kaizer/gologger/tags' -H 'authority: gitlab.sarafann.com' -H 'cache-control: max-age=0' -H 'origin: https://gitlab.sarafann.com' -H 'upgrade-insecure-requests: 1' -H 'content-type: application/x-www-form-urlencoded' -H 'user-agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.117 YaBrowser/20.2.0.1145 Yowser/2.5 Safari/537.36' -H 'sec-fetch-user: ?1' -H 'accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9' -H 'sec-fetch-site: same-origin' -H 'sec-fetch-mode: navigate' -H 'referer: https://gitlab.sarafann.com/kaizer/gologger/tags/new' -H 'accept-encoding: gzip, deflate, br' -H 'accept-language: ru,en;q=0.9' -H 'cookie: sidebar_collapsed=true; auto_devops_settings_dismissed=true; popup_video_number=1; geo_version=2.1; cookie_agreement=1; remember_user_token=W1syXSwiJDJhJDEwJEs0eERhZHJUUHBDT005SEVtMVFpY3UiLCIxNTgwMjgxOTU3LjMwMTA4MiJd--8904496755ffb386324b9ee5e0f5f5f067f4ebc2; CAN_PLAY_AUCTIONS=true; active_tab_id=3751580978159560; geo={%22geo_id%22:3707%2C%22geo_name%22:%22%D0%90%D0%BD%D0%B3%D0%B0%D1%80%D1%81%D0%BA%22%2C%22translate%22:%22angarsk%22%2C%22prepositional_name%22:%22%D0%90%D0%BD%D0%B3%D0%B0%D1%80%D1%81%D0%BA%D0%B5%22%2C%22geo_type%22:3%2C%22name%22:%22%D0%A0%D0%BE%D1%81%D1%81%D0%B8%D1%8F%2C%20%D0%98%D1%80%D0%BA%D1%83%D1%82%D1%81%D0%BA%D0%B0%D1%8F%20%D0%BE%D0%B1%D0%BB.%2C%20%D0%90%D0%BD%D0%B3%D0%B0%D1%80%D1%81%D0%BA%2C%20189-%D0%B9%20%D0%BA%D0%B2%D0%B0%D1%80%D1%82%D0%B0%D0%BB%2C%2013%22%2C%22latitude%22:52.50949096679688%2C%22longitude%22:103.83004760742188}; _gitlab_session=95dad5dc324adfa3a47721eb6606f4d6; event_filter=all' --data 'utf8=%E2%9C%93&authenticity_token=FOdVVdwa5jPVaN6EzKkxsd25%2FCK5yDOSE5YT4A%2Bfwwk%2FKVw%2Bb06qca6e9gn3Mu6yc%2Bm7pcLTqWabzZqJNEm%2FRA%3D%3D&tag_name=$version&ref=master&message=Release+$version&release_description=$releaseMsg' --compressed