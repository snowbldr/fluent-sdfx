// In-place STL viewer: turns a <figure> with a known stem into a live
// Three.js canvas the reader can spin around.
//
// Three is loaded lazily — the import only fires the first time the user
// clicks "View in 3D" on any figure, so the docs stay light by default.

let threePromise = null
function loadThree() {
  if (!threePromise) {
    threePromise = Promise.all([
      import('three'),
      import('three/addons/loaders/STLLoader.js'),
      import('three/addons/controls/OrbitControls.js')
    ]).then(([THREE, { STLLoader }, { OrbitControls }]) => ({ THREE, STLLoader, OrbitControls }))
  }
  return threePromise
}

export function attachViewer(figure, stlUrl) {
  if (figure.dataset.viewerAttached) return
  figure.dataset.viewerAttached = '1'

  const button = document.createElement('button')
  button.type = 'button'
  button.className = 'stl-viewer-toggle'
  button.innerHTML = '<span aria-hidden="true">⌬</span> View in 3D'
  button.addEventListener('click', () => {
    button.disabled = true
    button.textContent = 'Loading…'
    activate(figure, stlUrl).catch(err => {
      console.error('STL viewer:', err)
      button.disabled = false
      button.textContent = 'View in 3D'
    })
  })
  figure.appendChild(button)
}

async function activate(figure, stlUrl) {
  const { THREE, STLLoader, OrbitControls } = await loadThree()

  // Replace the static <img> with a live canvas wrapper. Keep the
  // <figcaption> in place so the caption still describes the part.
  const img = figure.querySelector('img')
  const wrapper = document.createElement('div')
  wrapper.className = 'stl-viewer-stage'
  if (img) {
    img.replaceWith(wrapper)
  } else {
    figure.prepend(wrapper)
  }
  const toggle = figure.querySelector('.stl-viewer-toggle')
  if (toggle) toggle.remove()

  // Three.js scene
  const scene = new THREE.Scene()
  // Read the docs theme from the root attribute so the viewer matches.
  const isDark = document.documentElement.getAttribute('data-theme') !== 'light'
  scene.background = new THREE.Color(isDark ? 0x14151a : 0xf5f6f9)

  const w = wrapper.clientWidth || 600
  const h = Math.round(w * 0.7)
  wrapper.style.height = h + 'px'

  // FOV chosen to match f3d's default (30°) so the static screenshot and
  // the live viewer frame the part the same way.
  const camera = new THREE.PerspectiveCamera(30, w / h, 0.1, 5000)

  const renderer = new THREE.WebGLRenderer({ antialias: true })
  renderer.setSize(w, h)
  renderer.setPixelRatio(Math.min(window.devicePixelRatio || 1, 2))
  wrapper.appendChild(renderer.domElement)

  // Three-point lighting that reads in both themes.
  scene.add(new THREE.AmbientLight(0xffffff, isDark ? 0.5 : 0.65))
  const key = new THREE.DirectionalLight(0xffffff, 1.1)
  key.position.set(2, 3, 1.5)
  scene.add(key)
  const fill = new THREE.DirectionalLight(0xffffff, 0.45)
  fill.position.set(-2, 1, -1)
  scene.add(fill)

  // Load STL.
  const loader = new STLLoader()
  const geometry = await new Promise((resolve, reject) => {
    loader.load(stlUrl, resolve, undefined, reject)
  })
  geometry.computeVertexNormals()
  geometry.computeBoundingBox()
  const center = geometry.boundingBox.getCenter(new THREE.Vector3())
  geometry.translate(-center.x, -center.y, -center.z)

  const material = new THREE.MeshStandardMaterial({
    color: isDark ? 0x9aa0b2 : 0x707888,
    roughness: 0.55,
    metalness: 0.05
  })
  const mesh = new THREE.Mesh(geometry, material)
  // sdfx is Z-up; three.js default is Y-up. Rotate so the model's Z axis
  // aligns with world Y.
  mesh.rotation.x = -Math.PI / 2
  scene.add(mesh)

  // Frame the part the same way f3d does: elevation 25°, azimuth 40°,
  // distance chosen to fit the bounding box's projection in the view plane.
  const elev = (25 * Math.PI) / 180
  const az = (40 * Math.PI) / 180
  const viewDir = new THREE.Vector3(
    Math.cos(elev) * Math.sin(az),
    Math.sin(elev),
    Math.cos(elev) * Math.cos(az)
  )

  // World-space bbox after the mesh rotation.
  mesh.updateMatrixWorld()
  const worldBox = new THREE.Box3().setFromObject(mesh)
  const wsz = worldBox.getSize(new THREE.Vector3())
  const halfFovV = (camera.fov * Math.PI) / 360
  // Compute the bbox's projected extent on the camera's image plane.
  const cameraUp = new THREE.Vector3(0, 1, 0)
  const cameraRight = new THREE.Vector3().crossVectors(cameraUp, viewDir).normalize()
  const cameraUpAdj = new THREE.Vector3().crossVectors(viewDir, cameraRight).normalize()
  let maxH = 0
  let maxV = 0
  for (let dx = -1; dx <= 1; dx += 2)
    for (let dy = -1; dy <= 1; dy += 2)
      for (let dz = -1; dz <= 1; dz += 2) {
        const c = new THREE.Vector3((dx * wsz.x) / 2, (dy * wsz.y) / 2, (dz * wsz.z) / 2)
        maxH = Math.max(maxH, Math.abs(c.dot(cameraRight)))
        maxV = Math.max(maxV, Math.abs(c.dot(cameraUpAdj)))
      }
  const aspect = w / h
  const distV = maxV / Math.tan(halfFovV)
  const distH = maxH / (Math.tan(halfFovV) * aspect)
  // 1.05 = ~5% padding around the part — matches f3d's default framing.
  const dist = Math.max(distV, distH) * 1.05

  camera.position.copy(viewDir.clone().multiplyScalar(dist))
  camera.lookAt(0, 0, 0)

  const controls = new OrbitControls(camera, renderer.domElement)
  controls.target.set(0, 0, 0)
  controls.enableDamping = true
  controls.dampingFactor = 0.08
  controls.update()

  // Keep canvas square-ish on resize.
  const resize = () => {
    const w = wrapper.clientWidth
    const h = Math.round(w * 0.7)
    wrapper.style.height = h + 'px'
    renderer.setSize(w, h)
    camera.aspect = w / h
    camera.updateProjectionMatrix()
  }
  const ro = new ResizeObserver(resize)
  ro.observe(wrapper)

  // React to theme toggles.
  const themeObs = new MutationObserver(() => {
    const dark = document.documentElement.getAttribute('data-theme') !== 'light'
    scene.background = new THREE.Color(dark ? 0x14151a : 0xf5f6f9)
    material.color.setHex(dark ? 0x9aa0b2 : 0x707888)
  })
  themeObs.observe(document.documentElement, { attributes: true, attributeFilter: ['data-theme'] })

  let frame
  const tick = () => {
    frame = requestAnimationFrame(tick)
    controls.update()
    renderer.render(scene, camera)
  }
  tick()

  // Stop the render loop if the viewer is removed from the page.
  const cleanup = new MutationObserver(() => {
    if (!document.body.contains(wrapper)) {
      cancelAnimationFrame(frame)
      ro.disconnect()
      themeObs.disconnect()
      cleanup.disconnect()
      renderer.dispose()
      geometry.dispose()
      material.dispose()
    }
  })
  cleanup.observe(document.body, { childList: true, subtree: true })
}
