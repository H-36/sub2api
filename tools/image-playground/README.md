# Image Playground Integration

Sub2API embeds `gpt_image_playground` as an isolated static sub-application.
The upstream snapshot lives in `third_party/gpt_image_playground/` and should
stay as close to upstream as possible.

## Layout

- `third_party/gpt_image_playground/`: upstream source snapshot.
- `tools/image-playground/build.sh`: builds the sub-app into Sub2API's frontend
  public directory.
- `frontend/public/image-playground-app/`: generated build output, ignored by
  git and copied into Sub2API's final frontend build.
- `frontend/src/views/user/ImagePlaygroundView.vue`: Sub2API iframe wrapper.

Sub2API-specific behavior is applied after the upstream build:

- PWA service worker and manifest links are removed.
- The iframe wrapper passes the current Sub2API theme and page background.
- The built image playground document receives a small theme bridge so its
  `html`, `body`, and `#root` backgrounds match Sub2API.
- Upstream's media-query dark CSS is converted to class-based dark CSS so the
  iframe follows Sub2API's theme toggle instead of the operating system theme.

## Build

From the repository root:

```bash
cd frontend
corepack pnpm build:image-playground
```

The normal frontend build runs this automatically before the Vue build.

## Sync Upstream

The current upstream revision is recorded in `tools/image-playground/UPSTREAM`.

Preferred update flow:

```bash
git subtree pull \
  --prefix=third_party/gpt_image_playground \
  https://github.com/CookSleep/gpt_image_playground.git \
  main \
  --squash
```

Keep Sub2API-specific integration changes outside `third_party/gpt_image_playground/`.
After syncing, rebuild:

```bash
cd frontend
corepack pnpm build:image-playground
```

If the subtree metadata is unavailable in a local branch, refresh the snapshot
with a clean upstream copy, then commit the vendor update separately from
Sub2API integration changes.
